//go:build wireinject
// +build wireinject

package main

import (
	"gateway/internal/application"

	"github.com/google/wire"

	"gateway/internal/client"
	"gateway/internal/config"
	"gateway/internal/facade"
	"gateway/internal/infras"
	"gateway/internal/server"
)

//go:generate go run github.com/google/wire/cmd/wire

func InitializeServer(cfg *config.Config) (*server.Server, error) {
	wire.Build(
		client.ProviderSet,
		application.ProviderSet,
		facade.ProviderSet,
		infras.ProviderSet,
		server.ProviderSet,
	)
	return nil, nil
}
