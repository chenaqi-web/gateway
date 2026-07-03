// Package grpc is the gRPC delivery layer (interface adapter).
// Handlers translate protobuf messages to application service calls and back.
package grpc

import (
	"net"

	"google.golang.org/grpc"

	"backend/core-server/internal/config"
)

// Server wraps the gRPC server lifecycle.
type Server struct {
	cfg  *config.Config
	grpc *grpc.Server
}

// New builds a gRPC server with all registered handlers.
func New(cfg *config.Config) (*Server, error) {
	s := grpc.NewServer()
	// register pb services here
	return &Server{cfg: cfg, grpc: s}, nil
}

// Start begins listening and serving gRPC requests.
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.cfg.Server.Addr)
	if err != nil {
		return err
	}
	return s.grpc.Serve(lis)
}

// Stop gracefully shuts down the gRPC server.
func (s *Server) Stop() {
	s.grpc.GracefulStop()
}
