package controller

import (
	"gateway/internal/application"
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
)

type HealthController struct{ svc *application.HealthService }

func NewHealthController(svc *application.HealthService) *HealthController {
	return &HealthController{svc: svc}
}

func (h *HealthController) Ping(c *gin.Context) {
	resp, err := h.svc.Ping(c.Request.Context())
	if err != nil {
		return
	}
	reponse.Success(c, resp)
}
