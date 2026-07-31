package controller

import (
	"backend/gateway/internal/client/rpc"
	"backend/gateway/internal/client/rpc/core-rpc/healthpb"
	"backend/gateway/internal/model/dto"
	"backend/gateway/internal/model/reponse"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
		reponse.Fail(c, http.StatusBadGateway, "core-server unavailable")
		return
	}

	reponse.Success(c, dto.PingResponse{Message: resp.GetMessage()})
}
