package models

// AllProviders 所有支持的提供商列表 (国内外共 20 家)
// 数据来源: latest_models.json (由 public-model-sync skill 维护)
var AllProviders = []ProviderInfo{
	// ---- 国际前 10 ----
	OpenAIProvider,
	AzureOpenAIProvider,
	AnthropicProvider,
	GoogleVertexProvider,
	BedrockProvider,
	MetaLlamaProvider,
	MistralProvider,
	CohereProvider,
	XAIProvider,
	PerplexityProvider,
	// ---- 国内前 10 ----
	QwenProvider,
	DoubaoProvider,
	ZhipuProvider,
	MoonshotProvider,
	DeepSeekProvider,
	MiniMaxProvider,
	BaiduProvider,
	HunyuanProvider,
	SparkProvider,
	Yi01AIProvider,
	// ---- 聚合 / 代理 ----
	OpenRouterProvider,
	OllamaProvider,
	LiteLLMProvider,
	MoyuAIProvider,
}

// GetProviderByID 根据 ID 获取提供商信息
func GetProviderByID(id string) *ProviderInfo {
	for i := range AllProviders {
		if AllProviders[i].ID == id {
			return &AllProviders[i]
		}
	}
	return nil
}

// GetModelByID 根据模型 ID 获取模型信息
func GetModelByID(modelID string) *ModelInfo {
	for i := range AllProviders {
		for j := range AllProviders[i].Models {
			if AllProviders[i].Models[j].ID == modelID {
				return &AllProviders[i].Models[j]
			}
		}
	}
	return nil
}

// GetAllModels 获取所有模型列表
func GetAllModels() []ModelInfo {
	var result []ModelInfo
	for _, p := range AllProviders {
		result = append(result, p.Models...)
	}
	return result
}

// GetProviderForModel 根据模型 ID 获取对应的提供商
func GetProviderForModel(modelID string) *ProviderInfo {
	for i := range AllProviders {
		for _, m := range AllProviders[i].Models {
			if m.ID == modelID {
				return &AllProviders[i]
			}
		}
	}
	return nil
}
