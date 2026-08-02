package client

import (
	"gateway/internal/client/http"
	"gateway/internal/client/rpc"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	rpc.NewRPCClient,
	http.NewHTTPClient,
)
