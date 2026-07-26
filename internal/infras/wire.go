package infras

import (
	"backend/gateway/internal/infras/api/llm"
	"backend/gateway/internal/infras/cache"
	"backend/gateway/internal/infras/clog"
	"backend/gateway/internal/infras/repo"
	"backend/gateway/internal/infras/storage"

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

var StorageProviderSet = wire.NewSet(
	storage.NewClient,
)

// ProviderSet 当前注入 cache / mysql / storage；clog 已就绪，后续需要时再并入
var ProviderSet = wire.NewSet(
	CacheProviderSet,
	RepoProviderSet,
	ApiProviderSet,
	StorageProviderSet,
)
