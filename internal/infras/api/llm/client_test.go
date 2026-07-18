package llm

import (
	"context"
	"strings"
	"testing"

	"backend/gateway/internal/config"
)

func newTestOllamaClient(t *testing.T) *Client {
	t.Helper()

	client, err := NewClient(&config.Config{
		LLM: config.LLMConfig{
			Provider:  "ollama",
			APIURL:    "http://127.0.0.1:11434",
			ModelID:   "llama3.2",
			ModelName: "Llama 3.2",
		},
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	return client
}

func TestOllamaChat(t *testing.T) {
	client := newTestOllamaClient(t)

	result, err := client.Chat(context.Background(), "你是一个简洁的助手，只回答一句话。", "1+1等于几？")
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	reply := strings.TrimSpace(result.Content)
	t.Logf("Chat reply: %s", reply)
	if reply == "" {
		t.Fatal("expected non-empty reply")
	}
}

func TestOllamaChatStream(t *testing.T) {
	client := newTestOllamaClient(t)

	var reply strings.Builder
	err := client.ChatStream(context.Background(), "你是一个简洁的助手，只回答一句话。", "北京是哪个国家的首都？", func(chunk string) error {
		reply.WriteString(chunk)
		return nil
	})
	if err != nil {
		t.Fatalf("ChatStream() error = %v", err)
	}

	content := strings.TrimSpace(reply.String())
	t.Logf("ChatStream reply: %s", content)
	if content == "" {
		t.Fatal("expected non-empty reply")
	}
}
