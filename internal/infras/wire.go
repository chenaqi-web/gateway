package infras

import (
	"gateway/internal/infras/api/llm"
	"gateway/internal/infras/cache"
	"gateway/internal/infras/clog"
	"gateway/internal/infras/repo"
	"gateway/internal/infras/storage"

	"github.com/google/wire"
)

var CacheProviderSet = wire.NewSet(
	cache.NewClient,
	cache.NewJwtBlacklist,
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
	LogProviderSet,
	CacheProviderSet,
	RepoProviderSet,
	ApiProviderSet,
	StorageProviderSet,
)
