package llm

import (
	"context"
	"fmt"
	"strings"

	"backend/gateway/internal/config"

	"github.com/tmc/langchaingo/llms"
)

type ChatResult struct {
	Content string
}

type StreamCallback func(chunk string) error

// provider 供应商接口
type provider interface {
	modelName() string
	chat(ctx context.Context, systemPrompt, userPrompt string) (*ChatResult, error)
	chatStream(ctx context.Context, systemPrompt, userPrompt string, callback StreamCallback) error
}

type Client struct {
	impl provider
}

func NewClient(cfg *config.Config) (*Client, error) {
	providerName := strings.ToLower(strings.TrimSpace(cfg.LLM.Provider))
	if providerName == "" {
		providerName = "ollama"
	}

	var impl provider
	var err error
	switch providerName {
	case "ollama":
		impl, err = newOllamaProvider(cfg)
	case "zhipu", "zhipuai":
		impl, err = newZhipuProvider(cfg)
	default:
		return nil, fmt.Errorf("unsupported llm provider: %s", providerName)
	}
	if err != nil {
		return nil, err
	}
	return &Client{impl: impl}, nil
}

func (c *Client) ModelName() string {
	return c.impl.modelName()
}

func (c *Client) Chat(ctx context.Context, systemPrompt, userPrompt string) (*ChatResult, error) {
	return c.impl.chat(ctx, systemPrompt, userPrompt)
}

func (c *Client) ChatStream(ctx context.Context, systemPrompt, userPrompt string, callback StreamCallback) error {
	return c.impl.chatStream(ctx, systemPrompt, userPrompt, callback)
}

func buildMessages(systemPrompt, userPrompt string) []llms.MessageContent {
	messages := make([]llms.MessageContent, 0, 2)
	if strings.TrimSpace(systemPrompt) != "" {
		messages = append(messages, llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt))
	}
	messages = append(messages, llms.TextParts(llms.ChatMessageTypeHuman, userPrompt))
	return messages
}
