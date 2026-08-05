package rpc

import (
	"fmt"
	"gateway/internal/client/rpc/core-rpc/articlepb"
	"gateway/internal/client/rpc/core-rpc/authpb"
	"gateway/internal/client/rpc/core-rpc/categorypb"
	"gateway/internal/client/rpc/core-rpc/commentpb"
	"gateway/internal/client/rpc/core-rpc/likepb"
	"gateway/internal/client/rpc/core-rpc/userpb"
	"time"

	"github.com/hashicorp/go-multierror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gateway/internal/client/rpc/core-rpc/healthpb"

	"gateway/internal/config"
)

type Client struct {
	// rpc 连接
	coreConnection *grpc.ClientConn

	// 不同的proto客户端
	healthClient   healthpb.HealthServiceClient
	UserClient     userpb.UserServiceClient
	CategoryClient categorypb.CategoryServiceClient
	ArticleClient  articlepb.ArticleServiceClient
	CommentClient  commentpb.CommentServiceClient
	LikeClient     likepb.LikeServiceClient
	authClient     authpb.AuthServiceClient

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
		CategoryClient: categorypb.NewCategoryServiceClient(coreConn),
		ArticleClient:  articlepb.NewArticleServiceClient(coreConn),
		CommentClient:  commentpb.NewCommentServiceClient(coreConn),
		LikeClient:     likepb.NewLikeServiceClient(coreConn),
		authClient:     authpb.NewAuthServiceClient(coreConn),
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

func (c *Client) GetAuthClient() authpb.AuthServiceClient {
	return c.authClient
}

//===============
