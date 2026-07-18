package application

import (
	"backend/gateway/internal/client/http"
	"backend/gateway/internal/config"
	"backend/gateway/internal/infras/api/llm"
	"backend/gateway/internal/infras/repo"
	"backend/gateway/internal/model/dto"
	"backend/gateway/internal/model/entity"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	defaultAiChatSessionTitle = "新建会话"
	defaultEmptyRAGReply      = "目前没有该相关知识"
	aiChatSystemPrompt        = `你的名字是「小爱」，是社团官网的ai助手`
)

var (
	ErrAiChatMissingContent   = errors.New("content is required")
	ErrAiChatMissingSessionID = errors.New("session_id is required")
	ErrAiChatSessionNotFound  = errors.New("session not found")
)

type AiChatService struct {
	cfg        *config.Config
	aiChatRepo *repo.AiChatRepo
	vector     *http.PyClient
	llm        *llm.Client
}

func NewAiChatService(
	cfg *config.Config,
	aiChatRepo *repo.AiChatRepo,
	vector *http.PyClient,
	llmClient *llm.Client,
) *AiChatService {
	return &AiChatService{
		cfg:        cfg,
		aiChatRepo: aiChatRepo,
		vector:     vector,
		llm:        llmClient,
	}
}

type AiChatStreamCallback func(chunk AiChatStreamChunk) error

type AiChatStreamChunk struct {
	SessionID string
	Content   string
	Done      bool
	Knowledge []string
}

func (s *AiChatService) CreateSession(ctx context.Context, userID string) (*entity.AiChatSession, error) {
	session := &entity.AiChatSession{
		UserID:    userID,
		SessionID: newSessionID(),
		Title:     defaultAiChatSessionTitle,
	}
	if err := s.aiChatRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *AiChatService) GetSession(ctx context.Context, userID, sessionID string) (*entity.AiChatSession, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, ErrAiChatMissingSessionID
	}
	session, err := s.aiChatRepo.GetSessionByUser(ctx, userID, sessionID)
	if err != nil {
		if errors.Is(err, repo.ErrAiChatSessionNotFound) {
			return nil, ErrAiChatSessionNotFound
		}
		return nil, err
	}
	return session, nil
}

func (s *AiChatService) ListSessions(ctx context.Context, userID string, page, pageSize int) ([]*entity.AiChatSession, error) {
	return s.aiChatRepo.ListSessionsByUser(ctx, userID, page, pageSize)
}

func (s *AiChatService) ListMessages(ctx context.Context, sessionID string) ([]*entity.AiChatMessage, error) {
	return s.aiChatRepo.ListMessagesBySession(ctx, sessionID)
}

func (s *AiChatService) Chat(ctx context.Context, userID, sessionID, content string, callback AiChatStreamCallback) error {
	// 1. 检索知识库
	knowledgeContents, err := s.searchKnowledge(ctx, content)
	if err != nil {
		return err
	}

	// 2. 拿到历史对话
	history, err := s.aiChatRepo.ListMessagesBySession(ctx, sessionID)
	if err != nil {
		return err
	}

	// 3. 构建系统提示词和用户提示词
	userPrompt, systemPrompt := buildAiChatPrompts(content, knowledgeContents, history)

	// 4. 发送给llm模型生成内容
	reply := strings.Builder{}
	emit := func(content string, done bool) error {
		if callback == nil {
			return nil
		}
		chunk := AiChatStreamChunk{
			SessionID: sessionID,
			Content:   content,
			Done:      done,
		}
		if done {
			chunk.Knowledge = knowledgeContents
		}
		return callback(chunk)
	}

	streamErr := s.llm.ChatStream(ctx, systemPrompt, userPrompt, func(chunk string) error {
		reply.WriteString(chunk)
		return emit(chunk, false)
	})
	if streamErr != nil {
		return streamErr
	}

	// 当模型返回空时
	finalReply := strings.TrimSpace(reply.String())
	if finalReply == "" {
		finalReply = defaultEmptyRAGReply
		if err := emit(finalReply, false); err != nil {
			return err
		}
		reply.Reset()
		reply.WriteString(finalReply)
	}

	if err := emit("", true); err != nil {
		return err
	}

	// 5.保存记录
	return s.saveChatRecords(ctx, userID, sessionID, content, reply.String())
}

// =====================================================================================================================

func (s *AiChatService) saveChatRecords(ctx context.Context, userID, sessionID, content, reply string) error {
	messages := []*entity.AiChatMessage{
		{
			SessionID: sessionID,
			UserID:    userID,
			Role:      entity.AiChatRoleUser,
			Content:   content,
		},
		{
			SessionID: sessionID,
			UserID:    userID,
			Role:      entity.AiChatRoleAssistant,
			Content:   reply,
		},
	}
	if err := s.aiChatRepo.CreateMessages(ctx, messages); err != nil {
		return err
	}
	return s.aiChatRepo.TouchSession(ctx, sessionID)
}

func (s *AiChatService) searchKnowledge(ctx context.Context, content string) ([]string, error) {
	collection := strings.TrimSpace(s.cfg.AiChat.VectorCollection)
	if collection == "" {
		collection = "test_knowledge"
	}
	topK := s.cfg.AiChat.DefaultTopK
	if topK <= 0 {
		topK = 5
	}

	resp, err := s.vector.SearchVectors(ctx, collection, &dto.VectorSearchRequest{
		Content: content,
		TopK:    &topK,
	})
	if err != nil {
		return nil, err
	}

	contents := make([]string, 0, len(resp.Data))
	for _, item := range resp.Data {
		knowledgeContent := strings.TrimSpace(item.Content)
		if knowledgeContent == "" {
			continue
		}
		contents = append(contents, knowledgeContent)
	}
	return contents, nil
}

func buildAiChatPrompts(question string, knowledgeContents []string, history []*entity.AiChatMessage) (userPrompt, systemPrompt string) {
	var systemBuilder strings.Builder
	systemBuilder.WriteString(aiChatSystemPrompt)

	// 1. 拼接RAG知识
	systemBuilder.WriteString("\n# RAG检索知识\n")
	if len(knowledgeContents) == 0 {
		systemBuilder.WriteString("（暂无匹配知识）\n")
	} else {
		for i, knowledgeContent := range knowledgeContents {
			systemBuilder.WriteString("[知识 kb-")
			systemBuilder.WriteString(strconv.Itoa(i + 1))
			systemBuilder.WriteString("] ")
			systemBuilder.WriteString(knowledgeContent)
			systemBuilder.WriteString("\n")
		}
	}
	systemBuilder.WriteString("\n")

	// 2. 拼接历史会话
	systemBuilder.WriteString("# 历史会话\n")
	if len(history) == 0 {
		systemBuilder.WriteString("（无）\n")
	} else {
		for _, record := range history {
			if record == nil || strings.TrimSpace(record.Content) == "" {
				continue
			}
			role := "你"
			if record.Role == entity.AiChatRoleUser {
				role = "用户"
			}
			systemBuilder.WriteString(role)
			systemBuilder.WriteString("：")
			systemBuilder.WriteString(record.Content)
			systemBuilder.WriteString("\n")
		}
	}

	// 3.构建用户提示词
	var userBuilder strings.Builder
	userBuilder.WriteString("# 用户问题\n")
	userBuilder.WriteString(question)
	userBuilder.WriteString("\n\n请根据以上RAG检索知识与历史会话，以小驭的身份生成可直接发送给用户的回复。")

	return userBuilder.String(), systemBuilder.String()
}

func newSessionID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}
