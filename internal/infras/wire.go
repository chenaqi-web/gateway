package infras

import (
	"backend/gateway/internal/infras/api/llm"
	"backend/gateway/internal/infras/cache"
	"backend/gateway/internal/infras/clog"
	"backend/gateway/internal/infras/repo"

	"github.com/google/wire"
)

var CacheProviderSet = wire.NewSet(
	cache.NewClient,
)

var RepoProviderSet = wire.NewSet(
	repo.NewDBClient,
	repo.NewAiChatRepo,
)

var LogProviderSet = wire.NewSet(
	clog.NewLog,
)

var ApiProviderSet = wire.NewSet(
	llm.NewClient,
)

// ProviderSet 当前注入 cache / mysql；clog 已就绪，后续需要时再并入
var ProviderSet = wire.NewSet(
	CacheProviderSet,
	RepoProviderSet,
	ApiProviderSet,
)
