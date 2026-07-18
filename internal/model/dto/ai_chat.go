package dto

import (
	"backend/gateway/internal/model/entity"
	"time"
)

type AiChatListSessionsQuery struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type AiChatChatRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// =====================================================================================================================

type AiChatSessionResponse struct {
	SessionID string `json:"session_id"`
	Title     string `json:"title"`
}

type AiChatMessageResponse struct {
	SessionID string    `json:"session_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type AiChatStreamChunkResponse struct {
	SessionID string   `json:"session_id"`
	Content   string   `json:"content"`
	Done      bool     `json:"done"`
	Knowledge []string `json:"knowledge,omitempty"`
}

func ToAiChatSessionResponse(session *entity.AiChatSession) *AiChatSessionResponse {
	if session == nil {
		return nil
	}
	return &AiChatSessionResponse{
		SessionID: session.SessionID,
		Title:     session.Title,
	}
}

func ToAiChatSessionResponses(sessions []*entity.AiChatSession) []*AiChatSessionResponse {
	list := make([]*AiChatSessionResponse, 0, len(sessions))
	for _, session := range sessions {
		list = append(list, ToAiChatSessionResponse(session))
	}
	return list
}

func ToAiChatMessageResponse(message *entity.AiChatMessage) *AiChatMessageResponse {
	if message == nil {
		return nil
	}
	return &AiChatMessageResponse{
		SessionID: message.SessionID,
		Role:      message.Role.String(),
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}
}

func ToAiChatMessageResponses(messages []*entity.AiChatMessage) []*AiChatMessageResponse {
	list := make([]*AiChatMessageResponse, 0, len(messages))
	for _, message := range messages {
		list = append(list, ToAiChatMessageResponse(message))
	}
	return list
}
