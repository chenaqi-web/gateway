package llm

import (
	"context"
	"fmt"
	"strings"

	"backend/gateway/internal/config"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type ollamaProvider struct {
	llm        *ollama.LLM
	ollamaName string
}

func newOllamaProvider(cfg *config.Config) (*ollamaProvider, error) {
	llmCfg := cfg.LLM
	serverURL := strings.TrimSpace(llmCfg.APIURL)
	if serverURL == "" {
		serverURL = "http://localhost:11434"
	}
	model := strings.TrimSpace(llmCfg.ModelID)
	if model == "" {
		model = "llama3.2"
	}
	modelName := strings.TrimSpace(llmCfg.ModelName)
	if modelName == "" {
		modelName = model
	}

	llm, err := ollama.New(
		ollama.WithModel(model),
		ollama.WithServerURL(serverURL),
	)
	if err != nil {
		return nil, fmt.Errorf("create ollama llm: %w", err)
	}

	return &ollamaProvider{
		llm:        llm,
		ollamaName: modelName,
	}, nil
}

func (p *ollamaProvider) modelName() string {
	return p.ollamaName
}

func (p *ollamaProvider) chat(ctx context.Context, systemPrompt, userPrompt string) (*ChatResult, error) {
	resp, err := p.llm.GenerateContent(ctx, buildMessages(systemPrompt, userPrompt))
	if err != nil {
		return nil, fmt.Errorf("ollama chat: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("ollama returned empty choices")
	}
	return &ChatResult{Content: resp.Choices[0].Content}, nil
}

func (p *ollamaProvider) chatStream(ctx context.Context, systemPrompt, userPrompt string, callback StreamCallback) error {
	_, err := p.llm.GenerateContent(ctx, buildMessages(systemPrompt, userPrompt), llms.WithStreamingFunc(
		func(_ context.Context, chunk []byte) error {
			if len(chunk) == 0 {
				return nil
			}
			return callback(string(chunk))
		},
	))
	if err != nil {
		return fmt.Errorf("ollama chat stream: %w", err)
	}
	return nil
}
