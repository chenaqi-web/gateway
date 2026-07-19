package rpc

import (
	"backend/gateway/internal/client/rpc/core-rpc/userpb"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"backend/gateway/internal/client/rpc/core-rpc/healthpb"

	"backend/gateway/internal/config"
)

type Client struct {
	// rpc 连接
	coreConnection *grpc.ClientConn

	// 不同的proto客户端
	healthClient healthpb.HealthServiceClient
	UserClient   userpb.UserServiceClient

	// 请求超时时间
	requestTimeout time.Duration
}

func NewRPCClient(cfg *config.Config) (*Client, error) {
	coreConn, err := newCoreConnection(cfg)
	if err != nil {
		return nil, err
	}

	timeoutSec := cfg.RPC.RequestTimeout
	if timeoutSec <= 0 {
		timeoutSec = 5
	}

	return &Client{
		coreConnection: coreConn,
		healthClient:   healthpb.NewHealthServiceClient(coreConn),
		UserClient:     userpb.NewUserServiceClient(coreConn),
		requestTimeout: time.Second * time.Duration(timeoutSec),
	}, nil
}

func newCoreConnection(cfg *config.Config) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		cfg.RPC.CoreServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial core-server %s: %w", cfg.RPC.CoreServerAddr, err)
	}
	return conn, nil
}

func (c *Client) Close() error {
	var multiErr *multierror.Error
	if err := c.coreConnection.Close(); err != nil {
		multiErr = multierror.Append(multiErr, err)
	}
	return multiErr.ErrorOrNil()
}

func (c *Client) GetRequestTimeout() time.Duration {
	return c.requestTimeout
}

// =====================================================================================================================
// 下面是健康检测

func (c *Client) GetHealthClient() healthpb.HealthServiceClient {
	return c.healthClient
}

func (c *Client) GetUserClient() userpb.UserServiceClient {
	return c.UserClient
}

//===============
