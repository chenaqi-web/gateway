package application

import (
	"context"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/userpb"
)

type UserService struct{ rpc *rpc.Client }

func NewUserService(rpcClient *rpc.Client) *UserService { return &UserService{rpc: rpcClient} }

func (s *UserService) Get(ctx context.Context) (any, error) {
	return s.rpc.UserClient.Login(ctx, &userpb.LoginReq{})
}
