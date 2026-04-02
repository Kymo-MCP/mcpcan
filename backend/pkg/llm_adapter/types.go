package llm_adapter

// llm_adapter/types.go
// 本文件直接 re-export pkg/llm 层的所有公共类型和常量，
// 避免重复定义，保证 adapter 层与底层实现的类型完全一致。
// ProviderType 常量、所有消息结构体（Message、ToolCall、ChatRequest 等）
// 均以类型别名方式暴露，调用方无需关心底层包名。

import (
	"context"
	"encoding/json"
	"os"

	orig_llm "github.com/kymo-mcp/mcpcan/pkg/llm"
)

// ——— Provider 类型 ———

// ProviderType 是 LLM 提供商的类型标识，直接复用底层定义
type ProviderType = orig_llm.ProviderType

// 所有支持的 Provider 常量，直接引用 pkg/llm 层
const (
	ProviderOpenAI       = orig_llm.ProviderOpenAI
	ProviderAzureOpenAI  = orig_llm.ProviderAzureOpenAI
	ProviderDeepSeek     = orig_llm.ProviderDeepSeek
	ProviderAnthropic    = orig_llm.ProviderAnthropic
	ProviderGoogle       = orig_llm.ProviderGoogle
	ProviderMistral      = orig_llm.ProviderMistral
	ProviderXAI          = orig_llm.ProviderXAI
	ProviderAzureBedrock = orig_llm.ProviderAzureBedrock
	ProviderVertexAI     = orig_llm.ProviderVertexAI
	ProviderMetaLlama    = orig_llm.ProviderMetaLlama
	ProviderCohere       = orig_llm.ProviderCohere
	ProviderPerplexity   = orig_llm.ProviderPerplexity
	// 聚合 / 代理
	ProviderOpenRouter = orig_llm.ProviderOpenRouter
	ProviderLiteLLM    = orig_llm.ProviderLiteLLM
	ProviderMoyuAI     = orig_llm.ProviderMoyuAI
	ProviderOllama     = orig_llm.ProviderOllama
	// 国内主流
	ProviderQwen     = orig_llm.ProviderQwen
	ProviderDoubao   = orig_llm.ProviderDoubao
	ProviderZhipu    = orig_llm.ProviderZhipu
	ProviderMoonshot = orig_llm.ProviderMoonshot
	ProviderBaidu    = orig_llm.ProviderBaidu
	ProviderHunyuan  = orig_llm.ProviderHunyuan
	ProviderSpark    = orig_llm.ProviderSpark
	ProviderMiniMax  = orig_llm.ProviderMiniMax
	ProviderYi01AI   = orig_llm.ProviderYi01AI
	// 特殊
	ProviderMCP = orig_llm.ProviderMCP
)

// SupportedProviders / DefaultBaseURLs / GetSupportedProviderList 均引用底层实现
var SupportedProviders = orig_llm.SupportedProviders
var DefaultBaseURLs = orig_llm.DefaultBaseURLs

// GetSupportedProviderList 返回所有支持的 provider ID 列表
func GetSupportedProviderList() []string {
	return orig_llm.GetSupportedProviderList()
}

// IsValidProvider 检查 provider 字符串是否有效
func IsValidProvider(provider string) bool {
	return orig_llm.IsValidProvider(provider)
}

// ——— 消息/请求/响应 类型别名（直接复用底层，zero copy）———

type (
	ProviderConfig     = orig_llm.ProviderConfig
	MessageContentPart = orig_llm.MessageContentPart
	MessageImageURL    = orig_llm.MessageImageURL
	Message            = orig_llm.Message
	ToolCall           = orig_llm.ToolCall
	ToolCallFunction   = orig_llm.ToolCallFunction
	Tool               = orig_llm.Tool
	Function           = orig_llm.Function
	ChatRequest        = orig_llm.ChatRequest
	ToolOutput         = orig_llm.ToolOutput
	StreamResponse     = orig_llm.StreamResponse
	Usage              = orig_llm.Usage
)

// ExtraContent 辅助类型（json 原始消息，与底层保持一致）
var _ json.RawMessage // 保证 import 不被 trim

// Provider 接口：直接复用底层定义
type Provider interface {
	StreamChat(ctx context.Context, req ChatRequest) (<-chan StreamResponse, error)
}

// GlobalProxyURL 全局代理配置
var GlobalProxyURL string

func init() {
	if v := os.Getenv("HTTP_PROXY"); v != "" {
		GlobalProxyURL = v
	}
}
