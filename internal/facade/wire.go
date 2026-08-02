package facade

import (
	"gateway/internal/facade/controller"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	// 注册处理器
	controller.NewUserController,
	controller.NewCategoryController,
	controller.NewArticleController,
	controller.NewHealthController,
	controller.NewAiChatController,
	controller.NewStorageController,
	controller.NewCommentController,
	New,
)
