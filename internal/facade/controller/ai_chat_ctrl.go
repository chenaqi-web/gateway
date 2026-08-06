package controller

import (
	"errors"
	"gateway/internal/application"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AiChatController struct{ svc *application.AiChatService }

func NewAiChatController(svc *application.AiChatService) *AiChatController {
	return &AiChatController{svc: svc}
}

func (ct *AiChatController) CreateSession(c *gin.Context) {
	userID, ok := currentAuthUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	session, err := ct.svc.CreateSession(c.Request.Context(), userID)
	if err != nil {
		aiChatError(c, err)
		return
	}
	reponse.Success(c, dto.ToAiChatSessionResponse(session))
}

func (ct *AiChatController) ListSessions(c *gin.Context) {
	var query dto.AiChatListSessionsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	userID, ok := currentAuthUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	list, err := ct.svc.ListSessions(c.Request.Context(), userID, query.Page, query.PageSize)
	if err != nil {
		aiChatError(c, err)
		return
	}
	reponse.Success(c, dto.ToAiChatSessionResponses(list))
}

func (ct *AiChatController) GetSession(c *gin.Context) {
	userID, ok := currentAuthUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	session, err := ct.svc.GetSession(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		aiChatError(c, err)
		return
	}
	reponse.Success(c, dto.ToAiChatSessionResponse(session))
}

func (ct *AiChatController) ListMessages(c *gin.Context) {
	list, err := ct.svc.ListMessages(c.Request.Context(), c.Param("id"))
	if err != nil {
		aiChatError(c, err)
		return
	}
	reponse.Success(c, dto.ToAiChatMessageResponses(list))
}

func (ct *AiChatController) Chat(c *gin.Context) {
	var input dto.AiChatChatRequest
	if !bindJSON(c, &input) {
		return
	}
	userID, ok := currentAuthUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	if err := ct.svc.Chat(c.Request.Context(), userID, input.SessionID, input.Content, func(chunk application.AiChatStreamChunk) error {
		c.SSEvent("message", dto.AiChatStreamChunkResponse{SessionID: chunk.SessionID, Content: chunk.Content, Done: chunk.Done, Knowledge: chunk.Knowledge})
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

func currentAuthUserID(c *gin.Context) (string, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return "", false
	}
	return strconv.FormatUint(userID, 10), true
}

func aiChatError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, application.ErrAiChatMissingContent), errors.Is(err, application.ErrAiChatMissingSessionID), errors.Is(err, application.ErrAiChatSessionNotFound):
		reponse.Fail(c, http.StatusBadRequest, err.Error())
	default:
		log.Printf("ai chat request failed: %v", err)
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
	}
}
