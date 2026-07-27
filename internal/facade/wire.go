package facade

import (
	"backend/gateway/internal/facade/controller"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	// 注册处理器
	controller.NewUserController,
	controller.NewHealthController,
	controller.NewAuthController,
	controller.NewAiChatController,
	controller.NewStorageController,
	New,
)
