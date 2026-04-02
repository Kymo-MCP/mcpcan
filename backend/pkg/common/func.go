package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/kymo-mcp/mcpcan/pkg/database/model"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetGatewayRoutePrefix() string {
	pathPrefix := os.Getenv(GatewayServerPrefix)
	if len(pathPrefix) == 0 {
		return GatewayRoutePrefix
	}
	return pathPrefix
}

func GetMarketRoutePrefix() string {
	pathPrefix := os.Getenv(MarketServerPrefix)
	if len(pathPrefix) == 0 {
		return MarketRoutePrefix
	}
	return pathPrefix
}

func GetAuthzRoutePrefix() string {
	pathPrefix := os.Getenv(AuthzServerPrefix)
	if len(pathPrefix) == 0 {
		return AuthzRoutePrefix
	}
	return pathPrefix
}

func GetMarketMcpHostingServersPrefix() string {
	return path.Join("/mcp-hosting", "servers")
}

func GetMarketMcpOpenServicePrefix() string {
	return path.Join(GetMarketRoutePrefix(), "services")
}

func SetKubeConfig(byteCfg []byte) *rest.Config {
	if len(byteCfg) == 0 {
		return nil
	}

	// First parse kubeconfig structure
	clientConfig, err := clientcmd.Load(byteCfg)
	if err != nil {
		fmt.Printf("Failed to load kubeconfig: %v\n", err)
		return nil
	}

	// Check and fix current-context issues
	if clientConfig.CurrentContext != "" {
		// Check if current-context exists in contexts
		if _, exists := clientConfig.Contexts[clientConfig.CurrentContext]; !exists {
			// If not exists, try to use the first available context
			for contextName := range clientConfig.Contexts {
				clientConfig.CurrentContext = contextName
				fmt.Printf("Fixed current-context from '%s' to '%s'\n", clientConfig.CurrentContext, contextName)
				break
			}
		}
	}

	// Handle backtick issues in server URL
	for clusterName, cluster := range clientConfig.Clusters {
		if strings.HasPrefix(cluster.Server, "`") && strings.HasSuffix(cluster.Server, "`") {
			cluster.Server = strings.Trim(cluster.Server, "`")
			clientConfig.Clusters[clusterName] = cluster
			fmt.Printf("Fixed server URL for cluster '%s': %s\n", clusterName, cluster.Server)
		}
	}

	// Re-serialize the fixed configuration
	fixedConfig, err := clientcmd.Write(*clientConfig)
	if err != nil {
		fmt.Printf("Failed to serialize fixed kubeconfig: %v\n", err)
		return nil
	}

	// Create REST config using the fixed configuration
	config, err := clientcmd.RESTConfigFromKubeConfig(fixedConfig)
	if err != nil {
		fmt.Printf("Failed to create REST config: %v\n", err)
		return nil
	}

	if config == nil {
		fmt.Printf("REST config is nil\n")
		return nil
	}

	return config
}

// createTargetProxyConfigForDefatuleHostingImg creates target proxy configuration
func CreateTargetProxyConfigForDefatuleHostingImg(serviceName string, servicePort int32, mcpName string, mcpProtocol model.McpProtocol) (*model.McpServersConfig, string) {
	addr := fmt.Sprintf("http://%s:%d", serviceName, servicePort)
	if mcpProtocol == model.McpProtocolSSE {
		addr += fmt.Sprintf("/%s", mcpProtocol.String())
	}
	if mcpProtocol == model.McpProtocolStreamableHttp {
		addr += fmt.Sprintf("/%s", "mcp")
	}
	return &model.McpServersConfig{
		McpServers: map[string]*model.McpConfig{
			mcpName: {
				Transport: mcpProtocol.String(),
				URL:       addr,
			},
		},
	}, addr
}

// createTargetProxyConfigForHttp creates target proxy configuration
func CreateTargetProxyConfigForHttp(serviceName string, servicePort int32, mcpName string, mcpProtocol model.McpProtocol, servicePath string) *model.McpServersConfig {
	addr := fmt.Sprintf("http://%s:%d%s", serviceName, servicePort, servicePath)
	return &model.McpServersConfig{
		McpServers: map[string]*model.McpConfig{
			mcpName: {
				Transport: mcpProtocol.String(),
				URL:       addr,
			},
		},
	}
}

// MarshalAndAssignConfig marshals and assigns config to json.RawMessage
func MarshalAndAssignConfig(config interface{}) (json.RawMessage, error) {
	b, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}

// GetImage returns image with mirror prefix or default registry (77kymo)
func GetImage(name string) string {
	mirror := os.Getenv("REGISTORY_IMAGE_MIRROR")
	if mirror == "" {
		mirror = "77kymo"
	}
	return fmt.Sprintf("%s/%s", mirror, name)
}

// mcp-hosting image 77kymo/mcp-hosting:latest
// GetMcpHostingImage returns mcp-hosting image
func GetMcpHostingImage() string {
	return GetImage("mcp-hosting:latest")
}

// mcp-sidecar image 77kymo/mcp-sidecar:latest
// GetSidecarImage returns mcp-sidecar image
func GetSidecarImage() string {
	return GetImage("mcp-sidecar:latest")
}

// openapi-to-mcp image default: 77kymo/openapi-to-mcp:v0.2.7
// GetOpenapiToMcpImage returns openapi-to-mcp image
func GetOpenapiToMcpImage() string {
	if image := strings.TrimSpace(os.Getenv("OPENAPI_TO_MCP_IMAGE")); image != "" {
		return image
	}
	return GetImage("openapi-to-mcp:v0.2.7")
}

// GetSidecarPort returns sidecar port from env or default 61180
func GetSidecarPort() int32 {
	port := os.Getenv(SidecarServerPortEnv)
	if port == "" {
		return 61180
	}
	var p int
	fmt.Sscanf(port, "%d", &p)
	if p == 0 {
		return 61180
	}
	return int32(p)
}

// GetMcpHostingPort returns mcp-hosting default port from env or default (depending on image/protocol)
// Note: Hosting port is usually passed from request, but this getter provides a fallback or default for construction
func GetMcpHostingPort() int32 {
	port := os.Getenv(HostingServerPortEnv)
	if port == "" {
		return 8080 // Default for most hosting images
	}
	var p int
	fmt.Sscanf(port, "%d", &p)
	if p == 0 {
		return 8080
	}
	return int32(p)
}
