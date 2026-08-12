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
	controller.NewAuthController,
	controller.NewAiChatController,
	controller.NewVectorController,
	controller.NewStorageController,
	controller.NewCommentController,
	controller.NewLikeController,
	New,
)
