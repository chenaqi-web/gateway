package controller

import (
	"backend/gateway/internal/application"
	"backend/gateway/internal/model/dto"
	"backend/gateway/internal/model/reponse"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AiChatController struct {
	svc *application.AiChatService
}

func NewAiChatController(svc *application.AiChatService) *AiChatController {
	return &AiChatController{svc: svc}
}

// =====================================================================================================================
// 会话方面的内容
// todo 后续需要加入用户id，目前还没有做中间件。

func (ct *AiChatController) CreateSession(c *gin.Context) {
	session, err := ct.svc.CreateSession(c.Request.Context(), "")
	if err != nil {
		aiChatError(c, err)
		return
	}
	c.JSON(http.StatusOK, reponse.Success(dto.ToAiChatSessionResponse(session)))
}

func (ct *AiChatController) ListSessions(c *gin.Context) {
	var query dto.AiChatListSessionsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, reponse.Error(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	list, err := ct.svc.ListSessions(c.Request.Context(), "", query.Page, query.PageSize)
	if err != nil {
		aiChatError(c, err)
		return
	}
	c.JSON(http.StatusOK, reponse.Success(dto.ToAiChatSessionResponses(list)))
}

func (ct *AiChatController) GetSession(c *gin.Context) {
	session, err := ct.svc.GetSession(c.Request.Context(), "", c.Param("id"))
	if err != nil {
		aiChatError(c, err)
		return
	}
	c.JSON(http.StatusOK, reponse.Success(dto.ToAiChatSessionResponse(session)))
}

// =====================================================================================================================
// chat方面

func (ct *AiChatController) ListMessages(c *gin.Context) {
	// url上的session-id拿到历史记录，唯一
	list, err := ct.svc.ListMessages(c.Request.Context(), c.Param("id"))
	if err != nil {
		aiChatError(c, err)
		return
	}
	c.JSON(http.StatusOK, reponse.Success(dto.ToAiChatMessageResponses(list)))
}

func (ct *AiChatController) Chat(c *gin.Context) {
	var input dto.AiChatChatRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, reponse.Error(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	// sse流式的头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	if err := ct.svc.Chat(c.Request.Context(), "", input.SessionID, input.Content, func(chunk application.AiChatStreamChunk) error {
		c.SSEvent("message", dto.AiChatStreamChunkResponse{
			SessionID: chunk.SessionID,
			Content:   chunk.Content,
			Done:      chunk.Done,
			Knowledge: chunk.Knowledge,
		})
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

func aiChatError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, application.ErrAiChatMissingContent),
		errors.Is(err, application.ErrAiChatMissingSessionID),
		errors.Is(err, application.ErrAiChatSessionNotFound):
		c.JSON(http.StatusBadRequest, reponse.Error(http.StatusBadRequest, err.Error()))
	default:
		log.Printf("ai chat request failed: %v", err)
		c.JSON(http.StatusInternalServerError, reponse.Error(http.StatusInternalServerError, "internal server error"))
	}
}
