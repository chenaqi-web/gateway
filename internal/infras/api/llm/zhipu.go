package llm

import (
	"context"
	"fmt"
	"strings"

	"backend/gateway/internal/config"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type zhipuProvider struct {
	llm       *openai.LLM
	ZhiPuName string
}

func newZhipuProvider(cfg *config.Config) (*zhipuProvider, error) {
	llmCfg := cfg.LLM
	apiKey := strings.TrimSpace(llmCfg.APIKey)
	if apiKey == "" {
		return nil, fmt.Errorf("zhipu api_key is required")
	}
	baseURL := strings.TrimSpace(llmCfg.APIURL)
	if baseURL == "" {
		baseURL = "https://open.bigmodel.cn/api/paas/v4/"
	}
	model := strings.TrimSpace(llmCfg.ModelID)
	if model == "" {
		model = "glm-4-flash"
	}
	modelName := strings.TrimSpace(llmCfg.ModelName)
	if modelName == "" {
		modelName = model
	}

	llm, err := openai.New(
		openai.WithToken(apiKey),
		openai.WithModel(model),
		openai.WithBaseURL(baseURL),
	)
	if err != nil {
		return nil, fmt.Errorf("create zhipu llm: %w", err)
	}

	return &zhipuProvider{
		llm:       llm,
		ZhiPuName: modelName,
	}, nil
}

func (p *zhipuProvider) modelName() string {
	return p.ZhiPuName
}

func (p *zhipuProvider) chat(ctx context.Context, systemPrompt, userPrompt string) (*ChatResult, error) {
	resp, err := p.llm.GenerateContent(ctx, buildMessages(systemPrompt, userPrompt))
	if err != nil {
		return nil, fmt.Errorf("zhipu chat: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("zhipu returned empty choices")
	}
	return &ChatResult{Content: resp.Choices[0].Content}, nil
}

func (p *zhipuProvider) chatStream(ctx context.Context, systemPrompt, userPrompt string, callback StreamCallback) error {
	_, err := p.llm.GenerateContent(ctx, buildMessages(systemPrompt, userPrompt), llms.WithStreamingFunc(
		func(_ context.Context, chunk []byte) error {
			if len(chunk) == 0 {
				return nil
			}
			return callback(string(chunk))
		},
	))
	if err != nil {
		return fmt.Errorf("zhipu chat stream: %w", err)
	}
	return nil
}
