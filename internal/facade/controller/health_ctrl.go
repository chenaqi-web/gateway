package controller

import (
	"context"
	"gateway/internal/application"
	"gateway/internal/model/reponse"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{ svc *application.HealthService }

func NewHealthController(svc *application.HealthService) *HealthController {
	return &HealthController{svc: svc}
}

func (h *HealthController) Ping(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5e9)
	defer cancel()
	resp, err := h.svc.Ping(ctx)
	if err != nil {
		log.Printf("health ping: %v", err)
		reponse.Fail(c, http.StatusBadGateway, "core-server unavailable")
		return
	}
	reponse.Success(c, resp)
}
