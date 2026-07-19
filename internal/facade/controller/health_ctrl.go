package controller

import (
	"backend/gateway/internal/model/dto"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/gateway/internal/client/rpc"
	"backend/gateway/internal/client/rpc/core-rpc/healthpb"
)

type HealthController struct {
	rpc *rpc.Client
}

func NewHealthController(rpcClient *rpc.Client) *HealthController {
	return &HealthController{rpc: rpcClient}
}

func (h *HealthController) Ping(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.rpc.GetRequestTimeout())
	defer cancel()

	resp, err := h.rpc.GetHealthClient().Ping(ctx, &healthpb.PingRequest{})
	if err != nil {
		log.Printf("health ping: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "core-server unavailable"})
		return
	}

	c.JSON(http.StatusOK, dto.PingResponse{Message: resp.GetMessage()})
}
