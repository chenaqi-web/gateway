package application

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAiChatService,
	NewVectorService,
	NewStorageService,
	NewArticleService,
	NewCategoryService,
	NewLikeService,
	NewCommentService,
	NewHealthService,
	NewUserService,
	NewAuthService,
)
