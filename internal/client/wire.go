package client

import (
	"backend/gateway/internal/client/rpc"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	rpc.NewRPCClient,
)
