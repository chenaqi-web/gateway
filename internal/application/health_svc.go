package application

import (
	"context"

	"gateway/internal/client/rpc"
	"gateway/internal/infras/clog"
	"gateway/internal/model/dto"
)

type HealthService struct {
	rpc *rpc.Client
	log *clog.Log
}

func NewHealthService(rpcClient *rpc.Client, log *clog.Log) *HealthService {
	return &HealthService{rpc: rpcClient, log: log}
}

func (s *HealthService) Ping(ctx context.Context) (*dto.PingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5e9)
	defer cancel()

	return &dto.PingResponse{Message: "success"}, nil
}
