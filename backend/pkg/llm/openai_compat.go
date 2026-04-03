package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// OpenAICompatProvider 是一个通用的 OpenAI 兼容 Provider 实现
// 直接使用 HTTP Client 调用 API，支持标准及扩展字段（如 reasoning_content）
// 替代 langchaingo 以获得对 API 字段的完整控制
type OpenAICompatProvider struct {
	BaseURL           string
	APIKey            string
	Client            *http.Client
	ExtraHeaders      map[string]string // 额外请求头（如 OpenRouter 的 X-Title）
	SupportsReasoning bool              // 是否支持 reasoning_content（Kimi/DeepSeek）
}

// NewOpenAICompatProvider creates a new OpenAICompatProvider
func NewOpenAICompatProvider(baseURL, apiKey string, client *http.Client, extraHeaders map[string]string, supportsReasoning bool) *OpenAICompatProvider {
	if client == nil {
		client = http.DefaultClient
	}
	return &OpenAICompatProvider{
		BaseURL:           baseURL,
		APIKey:            apiKey,
		Client:            client,
		ExtraHeaders:      extraHeaders,
		SupportsReasoning: supportsReasoning,
	}
}

// --- Request Types ---

// oaiMessage represents a message in OpenAI-compatible API format
type oaiMessage struct {
	Role             string        `json:"role"`
	Content          any           `json:"content"`
	ReasoningContent *string       `json:"reasoning_content,omitempty"` // nil=不序列化, ""=序列化为空字符串
	ToolCalls        []oaiToolCall `json:"tool_calls,omitempty"`
	ToolCallID       string        `json:"tool_call_id,omitempty"`
	Name             string        `json:"name,omitempty"`
}

type oaiToolCall struct {
	Index        *int            `json:"index,omitempty"`
	ID           string          `json:"id"`
	Type         string          `json:"type"`
	Function     oaiFunctionCall `json:"function"`
	ExtraContent json.RawMessage `json:"extra_content,omitempty"` // Google thought_signature 等
}

type oaiFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type oaiRequest struct {
	Model       string       `json:"model"`
	Messages    []oaiMessage `json:"messages"`
	Stream      bool         `json:"stream"`
	Temperature float32      `json:"temperature,omitempty"`
	Tools       []oaiTool    `json:"tools,omitempty"`
	ToolChoice  string       `json:"tool_choice,omitempty"`
}

type oaiTool struct {
	Type     string          `json:"type"`
	Function oaiToolFunction `json:"function"`
}

type oaiToolFunction struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Parameters  interface{} `json:"parameters,omitempty"`
}

// --- Response Types ---

type oaiStreamResponse struct {
	ID      string            `json:"id"`
	Choices []oaiStreamChoice `json:"choices"`
	Usage   *oaiUsage         `json:"usage,omitempty"`
}

type oaiUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type oaiStreamChoice struct {
	Index        int       `json:"index"`
	Delta        oaiDelta  `json:"delta"`
	FinishReason string    `json:"finish_reason"`
	Usage        *oaiUsage `json:"usage,omitempty"`
}

type oaiDelta struct {
	Role             string        `json:"role,omitempty"`
	Content          string        `json:"content,omitempty"`
	ReasoningContent string        `json:"reasoning_content,omitempty"`
	ToolCalls        []oaiToolCall `json:"tool_calls,omitempty"`
}

// getLocalImageBase64 检查相对路径，如果是本地 /static 图片则转为 base64
func getLocalImageBase64(urlPath string) string {
	if !strings.HasPrefix(urlPath, "/static/") {
		return urlPath
	}

	// 映射到本地物理路径
	// 当前写死为相对运行目录的 data/static。如果是正式服可通过 filepath.Join 和全局 config
	localPath := filepath.Join("data", urlPath)

	body, err := os.ReadFile(localPath)
	if err != nil {
		log.Printf("[OpenAICompat] Failed to read local image %s: %v", localPath, err)
		return urlPath
	}

	ext := strings.ToLower(filepath.Ext(localPath))
	mimeType := "image/jpeg"
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".webp":
		mimeType = "image/webp"
	case ".gif":
		mimeType = "image/gif"
	}

	b64 := base64.StdEncoding.EncodeToString(body)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, b64)
}

// StreamChat implements Provider interface
func (p *OpenAICompatProvider) StreamChat(ctx context.Context, req ChatRequest) (<-chan StreamResponse, error) {
	outCh := make(chan StreamResponse)

	// 1. Convert Messages
	var messages []oaiMessage
	for _, m := range req.Messages {
		var content any
		if len(m.MultiContent) > 0 {
			var multiParts []map[string]any
			for _, part := range m.MultiContent {
				p := map[string]any{"type": part.Type}
				if part.Type == "text" {
					p["text"] = part.Text
				} else if part.Type == "image_url" && part.ImageURL != nil {
					finalUrl := getLocalImageBase64(part.ImageURL.URL)
					p["image_url"] = map[string]string{"url": finalUrl}
				}
				multiParts = append(multiParts, p)
			}
			content = multiParts
		} else {
			content = m.Content
		}

		msg := oaiMessage{
			Role:       m.Role,
			Content:    content,
			ToolCallID: m.ToolCallID,
			Name:       m.ToolCallName,
		}

		// reasoning_content 处理逻辑：
		// - 如果有值，始终传递（*string 指向非空字符串）
		// - 如果无值且 SupportsReasoning=true 且是 assistant+tool_calls，传空字符串（强制序列化）
		// - 否则不传（nil，omitempty 跳过）
		if m.ReasoningContent != "" {
			rc := m.ReasoningContent
			msg.ReasoningContent = &rc
		} else if p.SupportsReasoning && m.Role == "assistant" && len(m.ToolCalls) > 0 {
			placeholder := " "
			msg.ReasoningContent = &placeholder
		}
		// 其他情况 msg.ReasoningContent 为 nil，omitempty 不序列化

		if len(m.ToolCalls) > 0 {
			for _, tc := range m.ToolCalls {
				msg.ToolCalls = append(msg.ToolCalls, oaiToolCall{
					ID:   tc.ID,
					Type: tc.Type,
					Function: oaiFunctionCall{
						Name:      tc.Function.Name,
						Arguments: tc.Function.Arguments,
					},
					ExtraContent: tc.ExtraContent, // 透传 Google thought_signature
				})
			}
		}
		messages = append(messages, msg)
	}

	// 2. Convert Tools
	var tools []oaiTool
	for _, t := range req.Tools {
		var params interface{} = t.Function.Parameters
		if raw, ok := params.(json.RawMessage); ok {
			var p map[string]interface{}
			if err := json.Unmarshal(raw, &p); err == nil {
				params = p
			}
		}
		tools = append(tools, oaiTool{
			Type: t.Type,
			Function: oaiToolFunction{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  params,
			},
		})
	}

	apiReq := oaiRequest{
		Model:       req.Model,
		Messages:    messages,
		Stream:      true,
		Temperature: req.Temperature,
		Tools:       tools,
	}

	// 部分 OpenAI 兼容网关在仅传 tools 时不会自动启用工具模式，显式声明 auto 提高兼容性
	if len(tools) > 0 {
		apiReq.ToolChoice = "auto"
	}

	reqBody, err := json.Marshal(apiReq)
	if err != nil {
		return nil, err
	}

	log.Printf("[OpenAICompat] POST %s/chat/completions model=%s msgs=%d tools=%d reasoning=%v", p.BaseURL, req.Model, len(messages), len(tools), p.SupportsReasoning)

	// 3. Create HTTP Request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.BaseURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)
	for k, v := range p.ExtraHeaders {
		httpReq.Header.Set(k, v)
	}

	// 4. Stream
	go func() {
		defer close(outCh)

		resp, err := p.Client.Do(httpReq)
		if err != nil {
			outCh <- StreamResponse{Error: err}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			outCh <- StreamResponse{Error: fmt.Errorf("API returned unexpected status code: %d: %s", resp.StatusCode, string(body))}
			return
		}

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					outCh <- StreamResponse{Error: err}
				}
				return
			}

			line = strings.TrimSpace(line)
			if line == "" || !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return
			}

			var streamResp oaiStreamResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}

			// Usage at top level
			if streamResp.Usage != nil {
				outCh <- StreamResponse{
					Usage: &Usage{
						PromptTokens:     streamResp.Usage.PromptTokens,
						CompletionTokens: streamResp.Usage.CompletionTokens,
						TotalTokens:      streamResp.Usage.TotalTokens,
					},
				}
			}

			if len(streamResp.Choices) > 0 {
				choice := streamResp.Choices[0]

				// Usage in choice
				if choice.Usage != nil {
					outCh <- StreamResponse{
						Usage: &Usage{
							PromptTokens:     choice.Usage.PromptTokens,
							CompletionTokens: choice.Usage.CompletionTokens,
							TotalTokens:      choice.Usage.TotalTokens,
						},
					}
				}

				sr := StreamResponse{}
				hasContent := false

				if choice.Delta.Content != "" {
					sr.Content = choice.Delta.Content
					hasContent = true
				}
				if choice.Delta.ReasoningContent != "" {
					sr.ReasoningContent = choice.Delta.ReasoningContent
					hasContent = true
				}
				if len(choice.Delta.ToolCalls) > 0 {
					for _, tc := range choice.Delta.ToolCalls {
						tcIndex := 0
						if tc.Index != nil {
							tcIndex = *tc.Index
						}
						sr.ToolCalls = append(sr.ToolCalls, ToolCall{
							Index: tcIndex,
							ID:    tc.ID,
							Type:  tc.Type,
							Function: ToolCallFunction{
								Name:      tc.Function.Name,
								Arguments: tc.Function.Arguments,
							},
							ExtraContent: tc.ExtraContent, // 捕获 Google thought_signature
						})
					}
					hasContent = true
				}

				if hasContent {
					select {
					case outCh <- sr:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return outCh, nil
}
