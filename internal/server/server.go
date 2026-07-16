package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"backend/gateway/internal/client/rpc"
	"backend/gateway/internal/config"
	"backend/gateway/internal/infras/cache"
	"backend/gateway/internal/infras/repo"
)

type Server struct {
	cfg         *config.Config
	Engine      *gin.Engine
	httpServer  *http.Server
	rpcClient   *rpc.Client
	dbClient    *repo.DBClient
	cacheClient *cache.CacheClient
}

func NewServer(
	cfg *config.Config,
	rpcClient *rpc.Client,
	dbClient *repo.DBClient,
	cacheClient *cache.CacheClient,
	engine *gin.Engine,
) (*Server, error) {
	return &Server{
		cfg:         cfg,
		Engine:      engine,
		rpcClient:   rpcClient,
		dbClient:    dbClient,
		cacheClient: cacheClient,
		httpServer: &http.Server{
			Addr:    cfg.Server.Addr,
			Handler: engine,
		},
	}, nil
}

func (s *Server) Start() error {
	log.Printf("gateway listening on %s", s.cfg.Server.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen %s: %w", s.cfg.Server.Addr, err)
	}
	return nil
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("shutdown server: %v", err)
	}
	if s.rpcClient != nil {
		if err := s.rpcClient.Close(); err != nil {
			log.Printf("close rpc client: %v", err)
		}
	}
	if s.dbClient != nil {
		if err := s.dbClient.Close(); err != nil {
			log.Printf("close mysql: %v", err)
		}
	}
	if s.cacheClient != nil {
		if err := s.cacheClient.Close(); err != nil {
			log.Printf("close redis: %v", err)
		}
	}
}

func (s *Server) Addr() string {
	return s.cfg.Server.Addr
}
