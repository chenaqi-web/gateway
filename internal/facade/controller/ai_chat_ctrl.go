package controller

import (
	"errors"
	"gateway/internal/application"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AiChatController struct{ svc *application.AiChatService }

func NewAiChatController(svc *application.AiChatService) *AiChatController {
	return &AiChatController{svc: svc}
}

// =====================================================================================================================
// 会话方面

func (ct *AiChatController) CreateSession(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	session, err := ct.svc.CreateSession(c.Request.Context(), userID)
	if err != nil {
		aiChatError(c, err)
		return
	}
	reponse.Success(c, session)
}

func (ct *AiChatController) ListSessions(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	list, err := ct.svc.ListSessions(c.Request.Context(), userID)
	if err != nil {
		aiChatError(c, err)
		return
	}
	reponse.Success(c, list)
}

func (ct *AiChatController) ListMessages(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	list, err := ct.svc.ListMessages(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		aiChatError(c, err)
		return
	}
	reponse.Success(c, list)
}

func (ct *AiChatController) DeleteSession(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	if err := ct.svc.DeleteSession(c.Request.Context(), userID, c.Param("id")); err != nil {
		aiChatError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (ct *AiChatController) UpdateSession(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	var input dto.AiChatUpdateSessionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	session, err := ct.svc.UpdateSession(c.Request.Context(), userID, c.Param("id"), input.Title)
	if err != nil {
		aiChatError(c, err)
		return
	}
	reponse.Success(c, session)
}

// =====================================================================================================================
// chat

func (ct *AiChatController) Chat(c *gin.Context) {
	var input dto.AiChatChatRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	if err := ct.svc.Chat(c.Request.Context(), userID, input.SessionID, input.Content, input.CollectionName, func(chunk dto.AiChatStreamChunkResponse) error {
		c.SSEvent("message", chunk)
		c.Writer.Flush()
		return nil
	}); err != nil {
		log.Printf("ai chat stream failed: %v", err)
		c.SSEvent("message", gin.H{"error": err.Error(), "done": true})
		c.Writer.Flush()
		return
	}
	c.SSEvent("message", "[DONE]")
	c.Writer.Flush()
}

// =====================================================================================================================
// settings

func (ct *AiChatController) GetStatus(c *gin.Context) {
	reponse.Success(c, ct.svc.GetStatus(c.Request.Context()))
}

func (ct *AiChatController) UpdateStatus(c *gin.Context) {
	if middleware.GetRole(c) != "admin" {
		reponse.Fail(c, http.StatusForbidden, "admin access required")
		return
	}
	var input dto.UpdateAiChatStatusRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	reponse.Success(c, ct.svc.UpdateStatus(c.Request.Context(), input.AssistantEnabled))
}

func (ct *AiChatController) GetSettings(c *gin.Context) {
	if middleware.GetRole(c) != "admin" {
		reponse.Fail(c, http.StatusForbidden, "admin access required")
		return
	}
	settings, err := ct.svc.GetSettings(c.Request.Context())
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, settings)
}

func (ct *AiChatController) UpdateSettings(c *gin.Context) {
	if middleware.GetRole(c) != "admin" {
		reponse.Fail(c, http.StatusForbidden, "admin access required")
		return
	}
	var input dto.UpdateAiSettingsRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	settings, err := ct.svc.UpdateSettings(c.Request.Context(), input)
	if err != nil {
		reponse.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	reponse.Success(c, settings)
}

func aiChatError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, application.ErrAiChatMissingContent), errors.Is(err, application.ErrAiChatMissingSessionID), errors.Is(err, application.ErrAiChatMissingSessionTitle), errors.Is(err, application.ErrAiChatSessionNotFound), errors.Is(err, application.ErrAiChatDisabled):
		reponse.Fail(c, http.StatusBadRequest, err.Error())
	default:
		log.Printf("ai chat request failed: %v", err)
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
	}
}
