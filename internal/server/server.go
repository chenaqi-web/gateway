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
	"backend/gateway/internal/facade/controller"
	"backend/gateway/internal/facade/router"
)

type Server struct {
	cfg        *config.Config
	Engine     *gin.Engine
	httpServer *http.Server
	rpcClient  *rpc.Client
}

func NewServer(cfg *config.Config, rpcClient *rpc.Client, health *controller.HealthController) (*Server, error) {
	engine := router.New(cfg, health)

	return &Server{
		cfg:       cfg,
		Engine:    engine,
		rpcClient: rpcClient,
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
		log.Printf("shutdown server server: %v", err)
	}
	if s.rpcClient != nil {
		if err := s.rpcClient.Close(); err != nil {
			log.Printf("close rpc client: %v", err)
		}
	}
}

func (s *Server) Addr() string {
	return s.cfg.Server.Addr
}
