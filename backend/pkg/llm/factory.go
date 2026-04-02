package llm

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/ollama"
)

// NewProvider creates a new LLM provider instance based on the type and config
func NewProvider(typ ProviderType, config ProviderConfig) (Provider, error) {
	var model llms.Model
	var err error
	// ctx := context.Background()

	// Create shared HTTP Client with Proxy if configured
	var baseTransport http.RoundTripper = http.DefaultTransport
	if config.ProxyURL != "" {
		proxyURL, err := url.Parse(config.ProxyURL)
		if err != nil {
			log.Printf("[Factory Warning] Invalid ProxyURL: %s, error: %v", config.ProxyURL, err)
		} else {
			baseTransport = &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					// Bypass proxy for local addresses
					host := req.URL.Hostname()
					if host == "localhost" || host == "127.0.0.1" || host == "::1" {
						return nil, nil
					}
					// Bypass private IP ranges (basic check)
					if strings.HasPrefix(host, "192.168.") || strings.HasPrefix(host, "10.") || strings.HasSuffix(host, ".local") {
						return nil, nil
					}
					// Check 172.16.x.x - 172.31.x.x
					if strings.HasPrefix(host, "172.") {
						parts := strings.Split(host, ".")
						if len(parts) >= 2 {
							if second, err := strconv.Atoi(parts[1]); err == nil {
								if second >= 16 && second <= 31 {
									return nil, nil
								}
							}
						}
					}
					return proxyURL, nil
				},
			}
			log.Printf("[Factory Info] Using Proxy: %s (ignoring local/private IPs)", config.ProxyURL)
		}
	}

	switch typ {
	case ProviderOpenAI, ProviderDeepSeek, ProviderMoonshot, ProviderQwen, ProviderDoubao, ProviderZhipu, ProviderXAI, ProviderMistral, ProviderOpenRouter, ProviderLiteLLM, ProviderMoyuAI, ProviderAzureOpenAI,
		// 新增国内厂商（都兼容 OpenAI Chat Completions 接口）
		ProviderBaidu, ProviderHunyuan, ProviderSpark, ProviderMiniMax, ProviderYi01AI,
		// 新增国际厂商（OpenAI 兼容接口）
		ProviderCohere, ProviderPerplexity:
		// 所有 OpenAI 兼容 Provider 统一使用自定义 HTTP 实现
		// 支持 reasoning_content 等标准扩展字段，不依赖 langchaingo
		baseURL := config.BaseURL
		if baseURL == "" {
			if defaultURL, ok := DefaultBaseURLs[typ]; ok && defaultURL != "" {
				baseURL = defaultURL
			}
		}

		httpClient := &http.Client{
			Transport: baseTransport,
		}

		var extraHeaders map[string]string
		if typ == ProviderOpenRouter {
			extraHeaders = map[string]string{
				"HTTP-Referer": "https://github.com/kymo-mcp/mcpcan",
				"X-Title":      "MCPCan",
			}
		}

		// reasoning_content 支持：基于 provider type 或 BaseURL 启用
		supportsReasoning := (typ == ProviderMoonshot || typ == ProviderDeepSeek ||
			strings.Contains(baseURL, "moonshot.cn") || strings.Contains(baseURL, "deepseek.com"))

		log.Printf("[Factory] Creating OpenAICompatProvider for %s, base URL: %s, reasoning=%v", typ, baseURL, supportsReasoning)
		return NewOpenAICompatProvider(baseURL, config.APIKey, httpClient, extraHeaders, supportsReasoning), nil

	// case ProviderGoogle:
	// 	// Google Gemini - Use OpenAI Compatibility Mode
	// 	// Docs: https://ai.google.dev/gemini-api/docs/openai
	// 	baseURL := config.BaseURL
	// 	if baseURL == "" {
	// 		baseURL = "https://generativelanguage.googleapis.com/v1beta/openai/"
	// 	}

	// 	httpClient := &http.Client{
	// 		Transport: baseTransport,
	// 	}

	// 	log.Printf("[Factory] Creating Google Provider (OpenAI Compat). BaseURL: %s", baseURL)
	// 	return NewOpenAICompatProvider(baseURL, config.APIKey, httpClient, nil, false), nil

	case ProviderAnthropic:
		// Anthropic Claude
		model, err = anthropic.New(
			anthropic.WithToken(config.APIKey),
			anthropic.WithHTTPClient(&http.Client{Transport: baseTransport}),
		)

	case ProviderOllama:
		// Ollama Local
		opts := []ollama.Option{}
		baseURL := config.BaseURL
		if baseURL != "" {
			baseURL = strings.TrimSuffix(baseURL, "/v1")
			baseURL = strings.TrimSuffix(baseURL, "/v1/")
			opts = append(opts, ollama.WithServerURL(baseURL))
		}
		// Ollama options usually include WithHTTPClient
		opts = append(opts, ollama.WithHTTPClient(&http.Client{Transport: baseTransport}))
		model, err = ollama.New(opts...)

	default:
		return nil, fmt.Errorf("unsupported provider type: %s", typ)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize provider %s: %w", typ, err)
	}

	return &LangChainAdapter{model: model}, nil
}

// headerTransport is a custom http.RoundTripper that adds headers to requests
type headerTransport struct {
	transport http.RoundTripper
	headers   map[string]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}

	// Debug logging
	// Avoid logging API Key
	logHeaders := make(http.Header)
	for k, v := range req.Header {
		if k == "Authorization" {
			logHeaders.Set(k, "Bearer ****")
		} else {
			logHeaders[k] = v
		}
	}
	log.Printf("[OpenRouter Debug] Request: %s %s, Headers: %v", req.Method, req.URL, logHeaders)

	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		// Read body for debugging
		bodyBytes, _ := io.ReadAll(resp.Body)
		// Restore body
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Printf("[OpenRouter Error] Status: %d, Body: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}
