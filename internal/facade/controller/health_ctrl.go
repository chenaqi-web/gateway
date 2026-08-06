package controller

import (
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
	resp, err := h.svc.Ping(c.Request.Context())
	if err != nil {
		log.Printf("health ping: %v", err)
		reponse.Fail(c, http.StatusBadGateway, "core-server unavailable")
		return
	}
	reponse.Success(c, resp)
}
