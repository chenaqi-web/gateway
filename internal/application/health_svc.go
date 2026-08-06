package application

import (
	"context"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/healthpb"
	"gateway/internal/model/dto"
)

type HealthService struct{ rpc *rpc.Client }

func NewHealthService(rpcClient *rpc.Client) *HealthService { return &HealthService{rpc: rpcClient} }

func (s *HealthService) Ping(ctx context.Context) (*dto.PingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5e9)
	defer cancel()

	resp, err := s.rpc.GetHealthClient().Ping(ctx, &healthpb.PingRequest{})
	if err != nil {
		return nil, err
	}
	return &dto.PingResponse{Message: resp.GetMessage()}, nil
}
