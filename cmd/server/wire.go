//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"backend/gateway/internal/client"
	"backend/gateway/internal/config"
	"backend/gateway/internal/facade"
	"backend/gateway/internal/server"
)

//go:generate go run github.com/google/wire/cmd/wire

func InitializeServer(cfg *config.Config) (*server.Server, error) {
	wire.Build(
		client.ProviderSet,
		facade.ProviderSet,
		server.ProviderSet,
	)
	return nil, nil
}
