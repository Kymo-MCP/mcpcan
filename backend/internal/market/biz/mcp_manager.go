package biz

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kymo-mcp/mcpcan/pkg/database/model"
	llm "github.com/kymo-mcp/mcpcan/pkg/llm_adapter"
	"github.com/kymo-mcp/mcpcan/pkg/utils"
	"github.com/kymo-mcp/mcpcan/pkg/version"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

type McpManager struct {
	clients      map[string]*client.Client
	configs      map[string]utils.McpServerConfig // 存储配置用于重连
	healthStatus map[string]bool                  // 存储健康状态
	mu           sync.RWMutex                     // 并发安全锁
	ctx          context.Context                  // 管理器生命周期上下文
	cancelFunc   context.CancelFunc               // 取消函数
}

func NewMcpManager() *McpManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &McpManager{
		clients:      make(map[string]*client.Client),
		configs:      make(map[string]utils.McpServerConfig),
		healthStatus: make(map[string]bool),
		ctx:          ctx,
		cancelFunc:   cancel,
	}
}

func (m *McpManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 取消上下文
	if m.cancelFunc != nil {
		m.cancelFunc()
	}

	// 关闭所有客户端
	for _, c := range m.clients {
		c.Close()
	}
}

func (m *McpManager) Initialize(ctx context.Context, configJSON string) error {
	if configJSON == "" || configJSON == "{}" {
		return nil
	}

	var config utils.McpServersConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("invalid mcp config: %v", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	successCount := 0
	for name, srv := range config.McpServers {
		// 容器内访问 MCP 网关时，将前端可见地址统一替换为内部入口地址
		// 仅对 /mcp-gateway/* 路径生效，避免影响 DIRECT 外部地址。
		srv.URL = normalizeGatewayURLForContainer(srv.URL)

		// 存储配置用于后续重连
		m.configs[name] = srv

		// 初始化客户端
		if err := m.initializeClient(ctx, name, srv); err != nil {
			fmt.Printf("[MCP Manager] failed to init client %s: %v\n", name, err)
			// 初始化失败,标记为不健康,但不返回错误
			m.healthStatus[name] = false
			continue
		}

		// 标记为健康
		m.healthStatus[name] = true
		successCount++
	}

	if len(config.McpServers) > 0 && successCount == 0 {
		return fmt.Errorf("all mcp servers failed to initialize")
	}

	return nil
}

func normalizeGatewayURLForContainer(rawURL string) string {
	internalBase := strings.TrimSpace(os.Getenv("MCP_INTERNAL_BASE_URL"))
	if rawURL == "" || internalBase == "" {
		return rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	// 仅处理 MCP 网关路由，避免影响直连外部 URL
	if !strings.HasPrefix(parsedURL.Path, "/mcp-gateway/") {
		return rawURL
	}

	internalURL, err := url.Parse(internalBase)
	if err != nil || internalURL.Scheme == "" || internalURL.Host == "" {
		return rawURL
	}

	parsedURL.Scheme = internalURL.Scheme
	parsedURL.Host = internalURL.Host
	return parsedURL.String()
}

// initializeClient 初始化单个 MCP 客户端 (内部方法,调用者需持有锁)
func (m *McpManager) initializeClient(ctx context.Context, name string, srv utils.McpServerConfig) error {
	var mcpClient *client.Client
	var err error

	// Determine transport type
	transportType := srv.Transport
	if transportType == "" {
		transportType = srv.Type
	}

	// Auto-detect if not specified
	if transportType == "" {
		if srv.URL != "" {
			if strings.HasSuffix(srv.URL, "/mcp") {
				transportType = model.McpProtocolStreamableHttp.String()
			} else {
				transportType = model.McpProtocolSSE.String()
			}
		} else if srv.Command != "" {
			transportType = model.McpProtocolStdio.String()
		}
	}

	transportType = strings.ToLower(transportType)
	transportType = strings.ReplaceAll(transportType, "_", "-") // Normalize snake_case to kebab-case if needed

	// 创建一个跳过证书验证的 HTTP Client 用于调试
	insecureClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	switch transportType {
	case model.McpProtocolSSE.String(), "http": // Standard MCP over HTTP (SSE)
		if srv.URL == "" {
			return fmt.Errorf("transport is %s but url is empty", transportType)
		}
		// Using NewSSEMCPClient for remote connections
		mcpClient, err = client.NewSSEMCPClient(
			srv.URL,
			client.WithHeaders(srv.Headers),
			client.WithHTTPClient(insecureClient),
		)
	case model.McpProtocolStreamableHttp.String(): // Streamable HTTP (NDJSON/etc)
		if srv.URL == "" {
			return fmt.Errorf("transport is %s but url is empty", transportType)
		}
		// Using NewStreamableHttpClient
		mcpClient, err = client.NewStreamableHttpClient(
			srv.URL,
			transport.WithHTTPHeaders(srv.Headers),
			transport.WithHTTPBasicClient(insecureClient),
		)
	case model.McpProtocolStdio.String():
		if srv.Command == "" {
			return fmt.Errorf("transport is stdio but command is empty")
		}
		// Using NewStdioMCPClient for local command execution
		var env []string
		for k, v := range srv.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		mcpClient, err = client.NewStdioMCPClient(
			srv.Command,
			env,
			srv.Args...,
		)
	default:
		return fmt.Errorf("unsupported transport: %s", transportType)
	}

	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	// Start client
	if err := mcpClient.Start(ctx); err != nil {
		return fmt.Errorf("failed to start client: %v", err)
	}

	_, err = mcpClient.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "mcpcan-backend",
				Version: version.Version,
			},
			Capabilities: mcp.ClientCapabilities{},
		},
	})
	if err != nil {
		mcpClient.Close()
		return fmt.Errorf("failed to initialize client: %v", err)
	}

	m.clients[name] = mcpClient
	return nil
}

// sanitizeServerName 将 mcp server name 的短横线替换为下划线，以符合 Google Tool Name 要求
// Google Tool Name 只能包含 letters, numbers, and underscores
func sanitizeServerName(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}

func (m *McpManager) GetTools(ctx context.Context) ([]llm.Tool, error) {
	var allTools []llm.Tool
	var listErrs []string

	for name, c := range m.clients {
		resp, err := c.ListTools(ctx, mcp.ListToolsRequest{})
		if err != nil {
			listErrs = append(listErrs, fmt.Sprintf("%s: %v", name, err))
			continue
		}

		safeName := sanitizeServerName(name)

		for _, t := range resp.Tools {
			toolName := fmt.Sprintf("%s__%s", safeName, t.Name)
			// 将 MCP ToolInputSchema 转为 map，确保包含 properties 字段
			// Google genai SDK 要求 type=object 的顶层 schema 必须包含 properties
			schemaBytes, _ := json.Marshal(t.InputSchema)
			var schemaMap map[string]interface{}
			if err := json.Unmarshal(schemaBytes, &schemaMap); err == nil {
				// 确保 properties 字段存在（Google 模型要求）
				props, hasProps := schemaMap["properties"].(map[string]interface{})
				if !hasProps {
					props = make(map[string]interface{})
					schemaMap["properties"] = props
				}

				// 如果 properties 为空，Google API 可能会拒绝（Error 400）。
				// Hack: 添加一个 _ignore 字段以满足 "properties map cannot be empty"
				if len(props) == 0 {
					props["_ignore"] = map[string]interface{}{
						"type":        "string",
						"description": "Ignored field to satisfy API requirements for empty parameters.",
					}
					schemaMap["properties"] = props // 更新回去
				}

				// 确保 type 字段存在
				if _, ok := schemaMap["type"]; !ok {
					schemaMap["type"] = "object"
				}
				schemaBytes, _ = json.Marshal(schemaMap)
			}

			llmTool := llm.Tool{
				Type: "function",
				Function: llm.Function{
					Name:        toolName,
					Description: t.Description,
					Parameters:  json.RawMessage(schemaBytes),
				},
			}
			allTools = append(allTools, llmTool)
		}
	}

	// 部分失败容错：只要还有可用工具就继续返回，避免单个 server 异常导致 tools 全量丢失
	if len(allTools) > 0 {
		if len(listErrs) > 0 {
			fmt.Printf("[MCP Manager] list tools partial success, errors=%s\n", strings.Join(listErrs, " | "))
		}
		return allTools, nil
	}

	if len(listErrs) > 0 {
		return nil, fmt.Errorf("failed to list tools from all mcp servers: %s", strings.Join(listErrs, " | "))
	}

	return allTools, nil
}

func (m *McpManager) CallTool(ctx context.Context, name string, args string) (*mcp.CallToolResult, error) {
	parts := strings.SplitN(name, "__", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid tool name format: %s", name)
	}
	sanitizedServerName := parts[0]
	toolName := parts[1]

	// 查找匹配的 client
	var c *client.Client

	m.mu.RLock()
	if client, ok := m.clients[sanitizedServerName]; ok {
		// 1. 如果正好名字没变（本来就没横线），直接命中
		c = client
	} else {
		// 2. 如果没直接命中，遍历查找哪个 original name 对应当前的 sanitized name
		for k, v := range m.clients {
			if sanitizeServerName(k) == sanitizedServerName {
				c = v
				break
			}
		}
	}
	m.mu.RUnlock()

	if c == nil {
		return nil, fmt.Errorf("server not found for sanitized name: %s", sanitizedServerName)
	}

	// 解析参数
	var argsMap map[string]interface{}
	if err := json.Unmarshal([]byte(args), &argsMap); err != nil {
		return nil, fmt.Errorf("invalid arguments json: %v", err)
	}

	// 如果有 _ignore 字段，移除它
	if _, ok := argsMap["_ignore"]; ok {
		delete(argsMap, "_ignore")
	}

	return c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: argsMap,
		},
	})
}
