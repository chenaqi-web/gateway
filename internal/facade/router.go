package facade

import (
	"gateway/internal/config"
	"gateway/internal/facade/controller"
	"gateway/internal/facade/middleware"
	"gateway/internal/facade/router"
	"gateway/internal/infras/cache"

	"github.com/gin-gonic/gin"
)

func New(cfg *config.Config,
	health *controller.HealthController,
	aiChat *controller.AiChatController,
	userCtrl *controller.UserController,
	categoryCtrl *controller.CategoryController,
	articleCtrl *controller.ArticleController,
	storageCtrl *controller.StorageController,
	commentCtrl *controller.CommentController,
	authCtrl *controller.AuthController,
	cacheClient *cache.CacheClient,
) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 跨域中间件
	r.Use(middleware.Cors())

	// 注册静态路由
	r.Static("/static", "./static")

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health/ping", health.Ping)
		router.NewAuthRouter(v1, authCtrl)

		v1.Use(middleware.RequireAuth(cfg.Auth, cacheClient))

		// 注册路由
		router.NewAIRouter(v1, aiChat)
		router.NewUserRouter(v1, userCtrl)
		router.NewCategoryRouter(v1, categoryCtrl)
		router.NewArticleRouter(v1, articleCtrl)
		router.NewStorageRouter(v1, storageCtrl)
		router.NewCommentRouter(v1, commentCtrl)

	}

	return r
}
