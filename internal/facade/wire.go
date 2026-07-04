package facade

import (
	"backend/gateway/internal/facade/controller"
	"backend/gateway/internal/facade/router"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	// 注册处理器
	controller.NewHealthController,

	// 注册路由
	router.New,
)
