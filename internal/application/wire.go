package application

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAiChatService,
	NewStorageService,
	NewArticleService,
	NewCategoryService,
	NewLikeService,
	NewCommentService,
	NewHealthService,
	NewUserService,
	NewAuthService,
)
