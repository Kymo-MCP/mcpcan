package llm

import (
	"context"
	"encoding/json"

	"github.com/kymo-mcp/mcpcan/pkg/llm/models"
)

// ProviderType defines the type of LLM provider
type ProviderType string

const (
	// 国际核心 Providers (OpenAI-compatible)
	ProviderOpenAI      ProviderType = "openai"
	ProviderAzureOpenAI ProviderType = "azure_openai"
	ProviderDeepSeek    ProviderType = "deepseek"
	ProviderAnthropic   ProviderType = "anthropic"
	ProviderGoogle      ProviderType = "google"
	ProviderMistral     ProviderType = "mistral"
	ProviderXAI         ProviderType = "xai"
	// 云厂商多模型网关 (国际)
	ProviderAzureBedrock ProviderType = "bedrock"
	ProviderVertexAI     ProviderType = "vertex_ai"
	ProviderMetaLlama    ProviderType = "meta_llama"
	ProviderCohere       ProviderType = "cohere"
	ProviderPerplexity   ProviderType = "perplexity"

	// 聚合 / 代理 Providers
	ProviderOpenRouter ProviderType = "openrouter"
	ProviderLiteLLM    ProviderType = "litellm"
	ProviderMoyuAI     ProviderType = "moyu_ai"
	ProviderOllama     ProviderType = "ollama"

	// 国内主流 Providers
	ProviderQwen     ProviderType = "qwen"     // 阿里云通义千问 (DashScope)
	ProviderDoubao   ProviderType = "doubao"   // 字节跳动豆包 (Volcengine)
	ProviderZhipu    ProviderType = "zhipu"    // 智谱 AI GLM
	ProviderMoonshot ProviderType = "moonshot" // 月之暗面 Kimi
	ProviderBaidu    ProviderType = "baidu"    // 百度文心 ERNIE
	ProviderHunyuan  ProviderType = "hunyuan"  // 腾讯混元
	ProviderSpark    ProviderType = "spark"    // 科大讯飞星火
	ProviderMiniMax  ProviderType = "minimax"  // MiniMax
	ProviderYi01AI   ProviderType = "yi_01ai"  // 零一万物

	// MCP 特殊 Provider
	ProviderMCP ProviderType = "mcp"
)

// SupportedProviders contains all valid provider types
var SupportedProviders = map[ProviderType]bool{
	ProviderOpenAI:       true,
	ProviderAzureOpenAI:  true,
	ProviderDeepSeek:     true,
	ProviderAnthropic:    true,
	ProviderGoogle:       true,
	ProviderMistral:      true,
	ProviderXAI:          true,
	ProviderAzureBedrock: true,
	ProviderVertexAI:     true,
	ProviderMetaLlama:    true,
	ProviderCohere:       true,
	ProviderPerplexity:   true,
	ProviderOpenRouter:   true,
	ProviderLiteLLM:      true,
	ProviderMoyuAI:       true,
	ProviderOllama:       true,
	ProviderQwen:         true,
	ProviderDoubao:       true,
	ProviderZhipu:        true,
	ProviderMoonshot:     true,
	ProviderBaidu:        true,
	ProviderHunyuan:      true,
	ProviderSpark:        true,
	ProviderMiniMax:      true,
	ProviderYi01AI:       true,
	ProviderMCP:          true,
}

// DefaultBaseURLs contains default API endpoints for each provider.
// Populated at init time from models.AllProviders to avoid duplicate maintenance.
// Fallback entries for providers not in AllProviders (e.g. MCP) are added manually.
var DefaultBaseURLs = map[ProviderType]string{
	ProviderMCP: "", // MCP 是本地协议，无需远端 BaseURL
}

func init() {
	// 从 models 包动态填充，消除与 providers_gen.go 中 BaseURL 的重复
	for _, p := range models.AllProviders {
		pt := ProviderType(p.ID)
		if _, exists := DefaultBaseURLs[pt]; !exists {
			DefaultBaseURLs[pt] = p.BaseURL
		}
	}
}

// GetSupportedProviderList returns list of all supported provider IDs
func GetSupportedProviderList() []string {
	return []string{
		string(ProviderOpenAI),
		string(ProviderAzureOpenAI),
		string(ProviderAnthropic),
		string(ProviderDeepSeek),
		string(ProviderGoogle),
		string(ProviderMistral),
		string(ProviderXAI),
		string(ProviderCohere),
		string(ProviderPerplexity),
		string(ProviderOpenRouter),
		string(ProviderLiteLLM),
		string(ProviderMoyuAI),
		string(ProviderOllama),
		string(ProviderQwen),
		string(ProviderDoubao),
		string(ProviderZhipu),
		string(ProviderMoonshot),
		string(ProviderBaidu),
		string(ProviderHunyuan),
		string(ProviderSpark),
		string(ProviderMiniMax),
		string(ProviderYi01AI),
		string(ProviderMCP),
	}
}

// IsValidProvider checks if the provider string is a valid provider type
func IsValidProvider(provider string) bool {
	return SupportedProviders[ProviderType(provider)]
}

// ProviderConfig holds configuration for LLM provider
type ProviderConfig struct {
	BaseURL  string
	APIKey   string
	ProxyURL string
}

// MessageContentPart represents a part of the message content (text or image)
type MessageContentPart struct {
	Type     string           `json:"type"` // text, image_url
	Text     string           `json:"text,omitempty"`
	ImageURL *MessageImageURL `json:"image_url,omitempty"`
}

// MessageImageURL represents image url
type MessageImageURL struct {
	URL string `json:"url"`
}

// Message represents a chat message
type Message struct {
	Role             string               `json:"role"`
	Content          string               `json:"content"`
	MultiContent     []MessageContentPart `json:"multi_content,omitempty"`
	ToolCalls        []ToolCall           `json:"tool_calls,omitempty"`
	ToolCallID       string               `json:"tool_call_id,omitempty"`
	ToolCallName     string               `json:"tool_call_name,omitempty"`    // Google FunctionResponse 需要函数名
	ReasoningContent string               `json:"reasoning_content,omitempty"` // 思考过程 (DeepSeek/Kimi)
}

// ToolCall represents a tool call request from LLM
type ToolCall struct {
	Index        int              `json:"index,omitempty"`
	ID           string           `json:"id"`
	Type         string           `json:"type"`
	Function     ToolCallFunction `json:"function"`
	ExtraContent json.RawMessage  `json:"extra_content,omitempty"` // 透传 Provider 扩展字段（如 Google thought_signature）
}

// ToolCallFunction represents the function call details in a tool call
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// Tool represents a tool definition
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function represents the function call details
type Function struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Parameters  interface{} `json:"parameters,omitempty"`
}

// ChatRequest represents a request to the LLM
type ChatRequest struct {
	Model       string
	Messages    []Message
	Tools       []Tool
	Temperature float32
	MaxTokens   int
	Stream      bool
}

// ToolOutput represents the result of a tool execution
type ToolOutput struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Result string `json:"result"`
}

// StreamResponse represents a chunk of response from streaming API
type StreamResponse struct {
	Content          string
	ReasoningContent string
	ToolCalls        []ToolCall
	ToolOutputs      []ToolOutput
	Usage            *Usage
	Error            error
}

// Usage represents token usage statistics
type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// Provider defines the interface for LLM providers
type Provider interface {
	StreamChat(ctx context.Context, req ChatRequest) (<-chan StreamResponse, error)
}
