// Code generated from latest_models.json; DO NOT EDIT manually.
// 使用 .agent/skills/public-model-sync/scripts/fetch_models.sh 重新生成后，手动更新此文件。
//
// 注意：附件能力字段（ImageMimeTypes/MaxImageSize/DocumentMimeTypes 等）
// 遵循 MCP 测试工具内部限制策略，比厂商官方更保守，见 model.go 中 MCP* 常量定义。
package models

// ================================
// 国际前 10 AI 厂商
// ================================

var OpenAIProvider = ProviderInfo{
	ID:          "openai",
	Name:        "OpenAI (GPT-5系列)",
	IconURL:     "https://cdn.openai.com/API/logo-dark.png",
	BaseURL:     "https://api.openai.com/v1",
	RegisterURL: "https://platform.openai.com/api-keys",
	DocsURL:     "https://platform.openai.com/docs",
	Models: []ModelInfo{
		{ID: "o4-mini", Name: "o4-mini", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: false, SupportTemperature: true, SupportThinking: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "openai"},
		{ID: "o4-mini-2025-04-16", Name: "o4-mini (2025-04-16)", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: false, SupportTemperature: true, SupportThinking: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "openai"},
		{ID: "o3", Name: "o3", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: false, SupportTemperature: false, SupportThinking: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "openai"},
		{ID: "o3-mini", Name: "o3-mini", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: false, SupportTemperature: false, SupportThinking: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "openai"},
		{ID: "o1", Name: "o1", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: false, SupportTemperature: false, SupportThinking: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "openai"},
		{ID: "gpt-4o", Name: "GPT-4o", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "openai"},
		{ID: "gpt-4o-mini", Name: "GPT-4o Mini", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "openai"},
		{ID: "gpt-4.1", Name: "GPT-4.1", ContextLength: 1047576, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "openai"},
		{ID: "gpt-4.1-mini", Name: "GPT-4.1 Mini", ContextLength: 1047576, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "openai"},
		{ID: "gpt-5", Name: "GPT-5", ContextLength: 1047576, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "openai"},
	},
}

var AzureOpenAIProvider = ProviderInfo{
	ID:          "azure_openai",
	Name:        "Microsoft Azure OpenAI",
	IconURL:     "https://upload.wikimedia.org/wikipedia/commons/4/44/Microsoft_logo.svg",
	BaseURL:     "",
	RegisterURL: "https://azure.microsoft.com/en-us/products/ai-services/openai-service",
	DocsURL:     "https://learn.microsoft.com/en-us/azure/ai-services/openai/",
	Models: []ModelInfo{
		{ID: "gpt-4o", Name: "GPT-4o", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "azure_openai"},
		{ID: "gpt-4o-mini", Name: "GPT-4o Mini", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "azure_openai"},
		{ID: "gpt-4.1", Name: "GPT-4.1", ContextLength: 272000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "azure_openai"},
		{ID: "gpt-5", Name: "GPT-5", ContextLength: 272000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "azure_openai"},
		{ID: "o3-mini", Name: "o3-mini", ContextLength: 200000, SupportTools: false, SupportSystemPrompt: false, SupportTemperature: false, SupportThinking: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "azure_openai"},
		{ID: "o1", Name: "o1", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: false, SupportTemperature: false, SupportThinking: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "azure_openai"},
	},
}

var AnthropicProvider = ProviderInfo{
	ID:          "anthropic",
	Name:        "Anthropic (Claude)",
	IconURL:     "https://www.anthropic.com/images/icons/safari-pinned-tab.svg",
	BaseURL:     "https://api.anthropic.com/v1",
	RegisterURL: "https://console.anthropic.com/",
	DocsURL:     "https://docs.anthropic.com/",
	Models: []ModelInfo{
		{ID: "claude-opus-4-6", Name: "Claude Opus 4.6", ContextLength: 1000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "anthropic"},
		{ID: "claude-opus-4-5", Name: "Claude Opus 4.5", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "anthropic"},
		{ID: "claude-sonnet-4-6", Name: "Claude Sonnet 4.6", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "anthropic"},
		{ID: "claude-sonnet-4-5", Name: "Claude Sonnet 4.5", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "anthropic"},
		{ID: "claude-sonnet-4-20250514", Name: "Claude Sonnet 4 (2025-05-14)", ContextLength: 1000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "anthropic"},
		{ID: "claude-opus-4-20250514", Name: "Claude Opus 4 (2025-05-14)", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "anthropic"},
		{ID: "claude-3-5-sonnet-20241022", Name: "Claude 3.5 Sonnet (Oct 2024)", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "anthropic"},
		{ID: "claude-3-5-haiku-20241022", Name: "Claude 3.5 Haiku", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, Provider: "anthropic"},
	},
}

var GoogleVertexProvider = ProviderInfo{
	ID:          "vertex_ai",
	Name:        "Google Vertex AI (Gemini)",
	IconURL:     "https://www.gstatic.com/lamda/images/gemini_sparkle_v002_d4735304ff6292a690345.svg",
	BaseURL:     "https://generativelanguage.googleapis.com/v1beta",
	RegisterURL: "https://aistudio.google.com/apikey",
	DocsURL:     "https://ai.google.dev/gemini-api/docs",
	Models: []ModelInfo{
		{ID: "gemini-2.5-pro", Name: "Gemini 2.5 Pro", ContextLength: 2000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "vertex_ai"},
		{ID: "gemini-2.5-flash", Name: "Gemini 2.5 Flash", ContextLength: 1048576, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, Provider: "vertex_ai"},
		{ID: "gemini-2.5-flash-lite", Name: "Gemini 2.5 Flash Lite", ContextLength: 1048576, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, Provider: "vertex_ai"},
		{ID: "gemini-2.0-flash", Name: "Gemini 2.0 Flash", ContextLength: 1048576, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, Provider: "vertex_ai"},
		{ID: "gemini-1.5-pro", Name: "Gemini 1.5 Pro", ContextLength: 2000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "vertex_ai"},
		{ID: "gemini-1.5-flash", Name: "Gemini 1.5 Flash", ContextLength: 1048576, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, Provider: "vertex_ai"},
	},
}

var BedrockProvider = ProviderInfo{
	ID:          "bedrock",
	Name:        "AWS Bedrock (多模型)",
	IconURL:     "https://a0.awsstatic.com/libra-css/images/logos/aws_logo_smile_1200x630.png",
	BaseURL:     "",
	RegisterURL: "https://aws.amazon.com/bedrock/",
	DocsURL:     "https://docs.aws.amazon.com/bedrock/",
	Models: []ModelInfo{
		{ID: "anthropic.claude-opus-4-20250514-v1:0", Name: "Claude Opus 4 (Bedrock)", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, Provider: "bedrock"},
		{ID: "anthropic.claude-sonnet-4-20250514-v1:0", Name: "Claude Sonnet 4 (Bedrock)", ContextLength: 1000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "bedrock"},
		{ID: "anthropic.claude-3-5-sonnet-20241022-v2:0", Name: "Claude 3.5 Sonnet v2 (Bedrock)", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "bedrock"},
		{ID: "meta.llama4-maverick-17b-instruct-v1:0", Name: "Llama 4 Maverick (Bedrock)", ContextLength: 1000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "bedrock"},
		{ID: "meta.llama3-1-405b-instruct-v1:0", Name: "Llama 3.1 405B (Bedrock)", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "bedrock"},
		{ID: "amazon.nova-pro-v1:0", Name: "Amazon Nova Pro", ContextLength: 300000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "bedrock"},
	},
}

var MetaLlamaProvider = ProviderInfo{
	ID:          "meta_llama",
	Name:        "Meta Llama 4",
	IconURL:     "https://upload.wikimedia.org/wikipedia/commons/thumb/7/7b/Meta_Platforms_Inc._logo.svg/800px-Meta_Platforms_Inc._logo.svg.png",
	BaseURL:     "https://api.llama.com/v1",
	RegisterURL: "https://llama.meta.com/",
	DocsURL:     "https://www.llama.com/docs/",
	Models: []ModelInfo{
		{ID: "Llama-4-Scout-17B-16E-Instruct-FP8", Name: "Llama 4 Scout 17B", ContextLength: 10000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "meta_llama"},
		{ID: "Llama-4-Maverick-17B-128E-Instruct-FP8", Name: "Llama 4 Maverick 17B", ContextLength: 1000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "meta_llama"},
		{ID: "Llama-3.3-70B-Instruct", Name: "Llama 3.3 70B", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "meta_llama"},
		{ID: "Llama-3.3-8B-Instruct", Name: "Llama 3.3 8B", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "meta_llama"},
	},
}

var MistralProvider = ProviderInfo{
	ID:          "mistral",
	Name:        "Mistral AI",
	IconURL:     "https://mistral.ai/images/mistral-logo-a9cd3e3d.svg",
	BaseURL:     "https://api.mistral.ai/v1",
	RegisterURL: "https://console.mistral.ai/",
	DocsURL:     "https://docs.mistral.ai/",
	Models: []ModelInfo{
		{ID: "mistral-large-latest", Name: "Mistral Large (Latest)", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "mistral"},
		{ID: "pixtral-large-latest", Name: "Pixtral Large (Latest)", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, Provider: "mistral"},
		{ID: "mistral-medium-3", Name: "Mistral Medium 3", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "mistral"},
		{ID: "mistral-small-latest", Name: "Mistral Small (Latest)", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "mistral"},
		{ID: "codestral-latest", Name: "Codestral (Latest)", ContextLength: 256000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "mistral"},
		{ID: "open-mixtral-8x22b", Name: "Mixtral 8x22B", ContextLength: 65336, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "mistral"},
	},
}

var CohereProvider = ProviderInfo{
	ID:          "cohere",
	Name:        "Cohere (Command R)",
	IconURL:     "https://cohere.com/favicon.ico",
	BaseURL:     "https://api.cohere.com/v1",
	RegisterURL: "https://dashboard.cohere.com/api-keys",
	DocsURL:     "https://docs.cohere.com/",
	Models: []ModelInfo{
		{ID: "command-a-03-2025", Name: "Command A (Mar 2025)", ContextLength: 256000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "cohere"},
		{ID: "command-r-plus", Name: "Command R+", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "cohere"},
		{ID: "command-r-plus-08-2024", Name: "Command R+ (Aug 2024)", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "cohere"},
		{ID: "command-r", Name: "Command R", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "cohere"},
		{ID: "command-r7b-12-2024", Name: "Command R7B (Dec 2024)", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "cohere"},
	},
}

var XAIProvider = ProviderInfo{
	ID:          "xai",
	Name:        "xAI (Grok)",
	IconURL:     "https://x.ai/favicon.ico",
	BaseURL:     "https://api.x.ai/v1",
	RegisterURL: "https://console.x.ai/",
	DocsURL:     "https://docs.x.ai/",
	Models: []ModelInfo{
		{ID: "grok-4", Name: "Grok 4", ContextLength: 256000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "xai"},
		{ID: "grok-4-latest", Name: "Grok 4 (Latest)", ContextLength: 256000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "xai"},
		{ID: "grok-3", Name: "Grok 3", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "xai"},
		{ID: "grok-3-fast", Name: "Grok 3 Fast", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "xai"},
		{ID: "grok-code-fast", Name: "Grok Code Fast", ContextLength: 256000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "xai"},
	},
}

var PerplexityProvider = ProviderInfo{
	ID:          "perplexity",
	Name:        "Perplexity AI (Sonar)",
	IconURL:     "https://www.perplexity.ai/favicon.ico",
	BaseURL:     "https://api.perplexity.ai",
	RegisterURL: "https://www.perplexity.ai/settings/api",
	DocsURL:     "https://docs.perplexity.ai/",
	Models: []ModelInfo{
		{ID: "sonar-pro", Name: "Sonar Pro", ContextLength: 200000, SupportTools: false, SupportSystemPrompt: true, SupportTemperature: true, Provider: "perplexity"},
		{ID: "sonar", Name: "Sonar", ContextLength: 128000, SupportTools: false, SupportSystemPrompt: true, SupportTemperature: true, Provider: "perplexity"},
		{ID: "sonar-reasoning-pro", Name: "Sonar Reasoning Pro", ContextLength: 128000, SupportTools: false, SupportSystemPrompt: false, SupportTemperature: true, SupportThinking: true, Provider: "perplexity"},
		{ID: "sonar-reasoning", Name: "Sonar Reasoning", ContextLength: 128000, SupportTools: false, SupportSystemPrompt: false, SupportTemperature: true, SupportThinking: true, Provider: "perplexity"},
		{ID: "sonar-deep-research", Name: "Sonar Deep Research", ContextLength: 128000, SupportTools: false, SupportSystemPrompt: false, SupportTemperature: true, Provider: "perplexity"},
	},
}

// ================================
// 国内前 10 AI 厂商
// ================================

var QwenProvider = ProviderInfo{
	ID:          "qwen",
	Name:        "阿里云通义千问 (Qwen)",
	IconURL:     "https://avatars.githubusercontent.com/u/148330874?s=200&v=4",
	BaseURL:     "https://dashscope.aliyuncs.com/compatible-mode/v1",
	RegisterURL: "https://bailian.console.aliyun.com/",
	DocsURL:     "https://help.aliyun.com/zh/model-studio/",
	Models: []ModelInfo{
		{ID: "qwen3-max", Name: "Qwen3 Max", ContextLength: 258048, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "qwen"},
		{ID: "qwen3-coder-plus", Name: "Qwen3 Coder Plus", ContextLength: 997952, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "qwen"},
		{ID: "qwq-plus", Name: "QwQ Plus", ContextLength: 98304, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportThinking: true, Provider: "qwen"},
		{ID: "qwen-turbo-latest", Name: "Qwen Turbo (Latest)", ContextLength: 1000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "qwen"},
		{ID: "qwen-plus-latest", Name: "Qwen Plus (Latest)", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "qwen"},
		{ID: "qwen-long", Name: "Qwen Long", ContextLength: 10000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "qwen"},
		{ID: "qwen2.5-72b-instruct", Name: "Qwen2.5 72B", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "qwen"},
		{ID: "qwen2.5-coder-32b-instruct", Name: "Qwen2.5 Coder 32B", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "qwen"},
	},
}

// DoubaoProvider 字节跳动豆包
// 官方图片限制：10MB/张，1张/次，格式：jpeg/png/webp/bmp/gif
// 官方文档限制：20MB/份，1个/次(会话最多3个)，格式：pdf/txt/csv/doc/docx/xls/xlsx/ppt/pptx/md
// MCP限制（更保守）：5MB，1个，仅 jpeg/png/pdf/txt
var DoubaoProvider = ProviderInfo{
	ID:          "doubao",
	Name:        "字节豆包 (Doubao)",
	IconURL:     "https://lf-flow-web-cdn.doubao.com/obj/flow-doubao/doubao/web/logo-icon.png",
	BaseURL:     "https://ark.cn-beijing.volces.com/api/v3",
	RegisterURL: "https://console.volcengine.com/ark",
	DocsURL:     "https://www.volcengine.com/docs/82379",
	Models: []ModelInfo{
		// Doubao Seed 2.0 系列（2026 旗舰，256K 上下文）
		{ID: "doubao-seed-2-0-pro-260215", Name: "豆包 Seed 2.0 Pro", ContextLength: 256000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: false, SupportThinking: true, SupportsVision: true, ImageMimeTypes: DoubaoImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "doubao"},
		{ID: "doubao-seed-2-0-lite", Name: "豆包 Seed 2.0 Lite", ContextLength: 256000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "doubao"},
		{ID: "doubao-seed-code-preview-251028", Name: "豆包 Seed 2.0 Code", ContextLength: 256000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "doubao"},
		// Doubao Seed 1.8 系列（2025-12 旗舰）
		{ID: "doubao-seed-1-8-251228", Name: "豆包 Seed 1.8", ContextLength: 256000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "doubao"},
		// Doubao 1.5 系列（2025-01）
		{ID: "doubao-1-5-pro-256k", Name: "豆包 1.5 Pro 256K", ContextLength: 256000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "doubao"},
		{ID: "doubao-1-5-pro-32k", Name: "豆包 1.5 Pro 32K", ContextLength: 32000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "doubao"},
		{ID: "doubao-1-5-lite-32k", Name: "豆包 1.5 Lite 32K", ContextLength: 32000, SupportTools: false, SupportSystemPrompt: true, SupportTemperature: true, Provider: "doubao"},
		{ID: "doubao-1-5-vision-pro-32k", Name: "豆包 1.5 Vision Pro 32K", ContextLength: 32000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: DoubaoImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, Provider: "doubao"},
		// Doubao 推理/思考模型
		{ID: "doubao-thinking-pro-250415", Name: "豆包 Thinking Pro", ContextLength: 32000, SupportTools: true, SupportSystemPrompt: false, SupportTemperature: false, SupportThinking: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "doubao"},
	},
}

var ZhipuProvider = ProviderInfo{
	ID:          "zhipu",
	Name:        "智谱 AI (GLM)",
	IconURL:     "https://www.zhipuai.cn/favicon.ico",
	BaseURL:     "https://open.bigmodel.cn/api/paas/v4",
	RegisterURL: "https://open.bigmodel.cn/",
	DocsURL:     "https://open.bigmodel.cn/dev/api",
	Models: []ModelInfo{
		{ID: "glm-4-plus", Name: "GLM-4 Plus", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsVision: true, ImageMimeTypes: CommonImageTypes, MaxImageSize: MCPMaxImageSize, MaxImageCount: MCPMaxImageCount, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "zhipu"},
		{ID: "glm-4-flash", Name: "GLM-4 Flash", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "zhipu"},
		{ID: "glm-4-long", Name: "GLM-4 Long", ContextLength: 1000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "zhipu"},
		{ID: "glm-4-air", Name: "GLM-4 Air", ContextLength: 8192, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "zhipu"},
		{ID: "glm-z1-air", Name: "GLM-Z1 Air", ContextLength: 32768, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportThinking: true, Provider: "zhipu"},
		{ID: "glm-z1-flash", Name: "GLM-Z1 Flash", ContextLength: 32768, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportThinking: true, Provider: "zhipu"},
	},
}

var MoonshotProvider = ProviderInfo{
	ID:          "moonshot",
	Name:        "月之暗面 (Kimi)",
	IconURL:     "https://statics.moonshot.cn/kimi-chat/favicon.ico",
	BaseURL:     "https://api.moonshot.cn/v1",
	RegisterURL: "https://platform.moonshot.cn/",
	DocsURL:     "https://platform.moonshot.cn/docs/api/chat",
	Models: []ModelInfo{
		{ID: "moonshot-v1-128k", Name: "Moonshot 128K", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "moonshot"},
		{ID: "moonshot-v1-auto", Name: "Moonshot Auto", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "moonshot"},
		{ID: "moonshot-v1-32k", Name: "Moonshot 32K", ContextLength: 32768, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "moonshot"},
		{ID: "moonshot-v1-8k", Name: "Moonshot 8K", ContextLength: 8192, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportsDocument: true, DocumentMimeTypes: BasicDocTypes, MaxDocumentSize: MCPMaxDocumentSize, MaxDocumentCount: MCPMaxDocumentCount, Provider: "moonshot"},
	},
}

var DeepSeekProvider = ProviderInfo{
	ID:          "deepseek",
	Name:        "DeepSeek",
	IconURL:     "https://chat.deepseek.com/favicon.ico",
	BaseURL:     "https://api.deepseek.com/v1",
	RegisterURL: "https://platform.deepseek.com/",
	DocsURL:     "https://api-docs.deepseek.com/",
	Models: []ModelInfo{
		{ID: "deepseek-v3", Name: "DeepSeek V3", ContextLength: 65536, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "deepseek"},
		{ID: "deepseek-v3.2", Name: "DeepSeek V3.2", ContextLength: 163840, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "deepseek"},
		{ID: "deepseek-chat", Name: "DeepSeek Chat", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "deepseek"},
		{ID: "deepseek-reasoner", Name: "DeepSeek R1", ContextLength: 131072, SupportTools: false, SupportSystemPrompt: true, SupportTemperature: true, SupportThinking: true, Provider: "deepseek"},
		{ID: "deepseek-r1", Name: "DeepSeek R1", ContextLength: 65536, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportThinking: true, Provider: "deepseek"},
		{ID: "deepseek-coder", Name: "DeepSeek Coder", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "deepseek"},
	},
}

var MiniMaxProvider = ProviderInfo{
	ID:          "minimax",
	Name:        "MiniMax (abab)",
	IconURL:     "https://www.minimaxi.com/favicon.ico",
	BaseURL:     "https://api.minimax.chat/v1",
	RegisterURL: "https://www.minimaxi.com/",
	DocsURL:     "https://platform.minimaxi.com/document/introduction",
	Models: []ModelInfo{
		{ID: "MiniMax-M2.5", Name: "MiniMax M2.5", ContextLength: 1000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "minimax"},
		{ID: "MiniMax-M2.5-lightning", Name: "MiniMax M2.5 Lightning", ContextLength: 1000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "minimax"},
		{ID: "MiniMax-M2", Name: "MiniMax M2", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "minimax"},
		{ID: "MiniMax-M2.1", Name: "MiniMax M2.1", ContextLength: 1000000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "minimax"},
	},
}

var BaiduProvider = ProviderInfo{
	ID:          "baidu",
	Name:        "百度文心 (ERNIE)",
	IconURL:     "https://nlp.bj.bcebos.com/ERNIE_logo.png",
	BaseURL:     "https://qianfan.baidubce.com/v2",
	RegisterURL: "https://console.bce.baidu.com/qianfan/ais/console/applicationConsole/application",
	DocsURL:     "https://cloud.baidu.com/doc/WENXINWORKSHOP/",
	Models: []ModelInfo{
		{ID: "ernie-4.5-turbo-128k", Name: "ERNIE 4.5 Turbo 128K", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "baidu"},
		{ID: "ernie-4.5-8k", Name: "ERNIE 4.5 8K", ContextLength: 8192, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "baidu"},
		{ID: "ernie-4.0-8k", Name: "ERNIE 4.0 8K", ContextLength: 8192, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "baidu"},
		{ID: "ernie-3.5-128k", Name: "ERNIE 3.5 128K", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "baidu"},
		{ID: "ernie-speed-128k", Name: "ERNIE Speed 128K", ContextLength: 131072, SupportTools: false, SupportSystemPrompt: true, SupportTemperature: true, Provider: "baidu"},
		{ID: "ernie-lite-8k", Name: "ERNIE Lite 8K", ContextLength: 8192, SupportTools: false, SupportSystemPrompt: true, SupportTemperature: true, Provider: "baidu"},
	},
}

var HunyuanProvider = ProviderInfo{
	ID:          "hunyuan",
	Name:        "腾讯混元 (Hunyuan)",
	IconURL:     "https://cloudcache.tencent-cloud.com/open_proj/proj_qcloud_v2/gateway/shareicons/cloud.png",
	BaseURL:     "https://api.hunyuan.cloud.tencent.com/v1",
	RegisterURL: "https://console.cloud.tencent.com/hunyuan/start",
	DocsURL:     "https://cloud.tencent.com/document/product/1729",
	Models: []ModelInfo{
		{ID: "hunyuan-turbo-latest", Name: "混元 Turbo (Latest)", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "hunyuan"},
		{ID: "hunyuan-pro", Name: "混元 Pro", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "hunyuan"},
		{ID: "hunyuan-standard", Name: "混元 Standard", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "hunyuan"},
		{ID: "hunyuan-code", Name: "混元 Code", ContextLength: 8192, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "hunyuan"},
		{ID: "hunyuan-lite", Name: "混元 Lite", ContextLength: 8192, SupportTools: false, SupportSystemPrompt: true, SupportTemperature: true, Provider: "hunyuan"},
	},
}

var SparkProvider = ProviderInfo{
	ID:          "spark",
	Name:        "科大讯飞星火 (Spark)",
	IconURL:     "https://xinghuo.xfyun.cn/favicon.ico",
	BaseURL:     "https://spark-api-open.xf-yun.com/v1",
	RegisterURL: "https://xinghuo.xfyun.cn/sparkapi",
	DocsURL:     "https://www.xfyun.cn/doc/spark/",
	Models: []ModelInfo{
		{ID: "spark-x1", Name: "星火 X1", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportThinking: true, Provider: "spark"},
		{ID: "spark-max", Name: "星火 Max", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "spark"},
		{ID: "spark-pro-128k", Name: "星火 Pro 128K", ContextLength: 131072, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "spark"},
		{ID: "spark-lite", Name: "星火 Lite", ContextLength: 4096, SupportTools: false, SupportSystemPrompt: true, SupportTemperature: true, Provider: "spark"},
	},
}

var Yi01AIProvider = ProviderInfo{
	ID:          "yi_01ai",
	Name:        "零一万物 (Yi)",
	IconURL:     "https://www.lingyiwanwu.com/favicon.ico",
	BaseURL:     "https://api.lingyiwanwu.com/v1",
	RegisterURL: "https://platform.lingyiwanwu.com/apikeys",
	DocsURL:     "https://platform.lingyiwanwu.com/docs",
	Models: []ModelInfo{
		{ID: "yi-lightning", Name: "Yi Lightning", ContextLength: 16384, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "yi_01ai"},
		{ID: "yi-large", Name: "Yi Large", ContextLength: 32768, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "yi_01ai"},
		{ID: "yi-large-turbo", Name: "Yi Large Turbo", ContextLength: 16384, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "yi_01ai"},
		{ID: "yi-medium", Name: "Yi Medium", ContextLength: 16384, SupportTools: false, SupportSystemPrompt: true, SupportTemperature: true, Provider: "yi_01ai"},
	},
}

// ================================
// 聚合 / 代理 Providers
// ================================

var OpenRouterProvider = ProviderInfo{
	ID:          "openrouter",
	Name:        "OpenRouter",
	IconURL:     "https://openrouter.ai/favicon.ico",
	BaseURL:     "https://openrouter.ai/api/v1",
	RegisterURL: "https://openrouter.ai/keys",
	DocsURL:     "https://openrouter.ai/docs",
	Models:      []ModelInfo{}, // 动态模型，通过 API 获取
}

var OllamaProvider = ProviderInfo{
	ID:          "ollama",
	Name:        "Ollama (本地)",
	IconURL:     "https://ollama.com/public/ollama.png",
	BaseURL:     "http://localhost:11434",
	RegisterURL: "https://ollama.com/",
	DocsURL:     "https://ollama.com/library",
	Models:      []ModelInfo{}, // 本地安装的模型，动态获取
}

var LiteLLMProvider = ProviderInfo{
	ID:          "litellm",
	Name:        "LiteLLM (代理)",
	IconURL:     "https://litellm.ai/favicon.ico",
	BaseURL:     "http://localhost:4000",
	RegisterURL: "https://litellm.vercel.app/",
	DocsURL:     "https://docs.litellm.ai/",
	Models:      []ModelInfo{}, // 通过 LiteLLM 代理的模型，动态获取
}

var MoyuAIProvider = ProviderInfo{
	ID:          "moyu_ai",
	Name:        "魔芋AI (New-API)",
	IconURL:     "https://www.moyu.info/logo.png?v=2&_t=493083",
	BaseURL:     "https://www.moyu.info/v1",
	RegisterURL: "https://www.moyu.info",
	DocsURL:     "http://101.200.167.88:8003/",
	Models: []ModelInfo{
		// 默认展示前 10 个模型（用于快速接入验证，实际可按账号权限调整）
		{ID: "chatgpt-4o-latest", Name: "ChatGPT-4o Latest", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "moyu_ai"},
		{ID: "gpt-4o", Name: "GPT-4o", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "moyu_ai"},
		{ID: "gpt-4o-mini", Name: "GPT-4o Mini", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "moyu_ai"},
		{ID: "gpt-4.1-mini", Name: "GPT-4.1 Mini", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "moyu_ai"},
		{ID: "claude-opus-4-6", Name: "Claude Opus 4.6", ContextLength: 200000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "moyu_ai"},
		{ID: "deepseek-v3", Name: "DeepSeek V3", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "moyu_ai"},
		{ID: "deepseek-r1", Name: "DeepSeek R1", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, SupportThinking: true, Provider: "moyu_ai"},
		{ID: "qwen-max", Name: "Qwen Max", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "moyu_ai"},
		{ID: "doubao-seedream-4-5-251128", Name: "Doubao Seedream 4.5", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "moyu_ai"},
		{ID: "doubao-seedance-1-5-pro-251215", Name: "Doubao Seedance 1.5 Pro", ContextLength: 128000, SupportTools: true, SupportSystemPrompt: true, SupportTemperature: true, Provider: "moyu_ai"},
	},
}

// GoogleProvider 保留向后兼容（映射到 vertex_ai）
var GoogleProvider = GoogleVertexProvider
