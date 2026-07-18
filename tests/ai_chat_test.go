package tests

import (
	"context"
	"os"
	"strings"
	"testing"

	"backend/gateway/internal/application"
	httpclient "backend/gateway/internal/client/http"
	"backend/gateway/internal/config"
	"backend/gateway/internal/infras/api/llm"
	"backend/gateway/internal/infras/repo"
	"backend/gateway/internal/model/entity"
)

func TestMain(m *testing.M) {
	if _, err := os.Stat("conf/config.yaml"); err != nil {
		if err := os.Chdir(".."); err != nil {
			panic(err)
		}
	}
	os.Exit(m.Run())
}

// 依赖本地服务：
// - MySQL（conf/config.yaml）
// - agent-server 向量检索（默认 http://127.0.0.1:8080）
// - Ollama（默认 http://127.0.0.1:11434，模型 llama3.2）
func newAiChatService(t *testing.T) *application.AiChatService {
	t.Helper()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	dbClient, err := repo.NewDBClient(cfg)
	if err != nil {
		t.Fatalf("connect mysql: %v", err)
	}
	t.Cleanup(func() {
		_ = dbClient.Close()
	})

	llmClient, err := llm.NewClient(cfg)
	if err != nil {
		t.Fatalf("create llm client: %v", err)
	}

	return application.NewAiChatService(
		cfg,
		repo.NewAiChatRepo(dbClient),
		httpclient.NewHTTPClient(cfg),
		llmClient,
	)
}

func TestAiChatServiceChat(t *testing.T) {
	svc := newAiChatService(t)
	ctx := context.Background()
	userID := "integration-test-user"

	session, err := svc.CreateSession(ctx, userID)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	t.Logf("created session: %s", session.SessionID)

	var chunks []application.AiChatStreamChunk
	reply := strings.Builder{}
	err = svc.Chat(ctx, userID, session.SessionID, "1+1等于几？", func(chunk application.AiChatStreamChunk) error {
		chunks = append(chunks, chunk)
		if !chunk.Done && chunk.Content != "" {
			reply.WriteString(chunk.Content)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	content := strings.TrimSpace(reply.String())
	t.Logf("reply: %s", content)
	if content == "" {
		t.Fatal("expected non-empty reply")
	}

	var doneChunk *application.AiChatStreamChunk
	for i := range chunks {
		if chunks[i].Done {
			doneChunk = &chunks[i]
			break
		}
	}
	if doneChunk == nil {
		t.Fatal("expected done chunk")
	}
	t.Logf("knowledge count: %d", len(doneChunk.Knowledge))

	messages, err := svc.ListMessages(ctx, session.SessionID)
	if err != nil {
		t.Fatalf("ListMessages() error = %v", err)
	}
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}
	if messages[0].Role != entity.AiChatRoleUser || messages[0].Content != "1+1等于几？" {
		t.Fatalf("unexpected user message: %+v", messages[0])
	}
	if messages[1].Role != entity.AiChatRoleAssistant || strings.TrimSpace(messages[1].Content) == "" {
		t.Fatalf("unexpected assistant message: %+v", messages[1])
	}
	t.Logf("saved assistant reply: %s", messages[1].Content)
}

func TestAiChatServiceChatWithHistory(t *testing.T) {
	svc := newAiChatService(t)
	ctx := context.Background()
	userID := "integration-test-user"

	session, err := svc.CreateSession(ctx, userID)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	firstQuestion := "你好"
	err = svc.Chat(ctx, userID, session.SessionID, firstQuestion, func(application.AiChatStreamChunk) error {
		return nil
	})
	if err != nil {
		t.Fatalf("first Chat() error = %v", err)
	}

	history, err := svc.ListMessages(ctx, session.SessionID)
	if err != nil {
		t.Fatalf("ListMessages() after first chat error = %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("expected 2 history messages after first chat, got %d", len(history))
	}

	secondQuestion := "我刚才问了什么？"
	reply := strings.Builder{}
	err = svc.Chat(ctx, userID, session.SessionID, secondQuestion, func(chunk application.AiChatStreamChunk) error {
		if !chunk.Done && chunk.Content != "" {
			reply.WriteString(chunk.Content)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("second Chat() error = %v", err)
	}

	content := strings.TrimSpace(reply.String())
	t.Logf("second reply: %s", content)
	if content == "" {
		t.Fatal("expected non-empty second reply")
	}

	messages, err := svc.ListMessages(ctx, session.SessionID)
	if err != nil {
		t.Fatalf("ListMessages() after second chat error = %v", err)
	}
	if len(messages) != 4 {
		t.Fatalf("expected 4 messages after second chat, got %d", len(messages))
	}
}
