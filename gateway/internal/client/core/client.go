// Package core is the gRPC client for core-server.
package core

import "backend/gateway/internal/config"

// Client wraps the gRPC connection to core-server.
type Client struct {
	addr string
}

// New creates a gRPC client pointed at core-server.
func New(cfg *config.GRPCConfig) (*Client, error) {
	return &Client{addr: cfg.CoreServerAddr}, nil
}

// Close releases the underlying gRPC connection.
func (c *Client) Close() error {
	return nil
}
