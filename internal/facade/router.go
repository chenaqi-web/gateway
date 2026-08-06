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
	likeCtrl *controller.LikeController,
	authCtrl *controller.AuthController,
	jwtBlackList *cache.JwtBlacklist,
) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 跨域中间件
	r.Use(middleware.Cors())

	// 认证中间件
	authMiddleware := middleware.NewAuthMiddleware(cfg.Auth, jwtBlackList)
	auth := authMiddleware.RequireAuth()

	// Prometheus metrics middleware
	//r.Use(infraProm.GinMiddleware())
	// r.Use(middleware.Recovery(), middleware.CORS())

	// 注册静态路由（含 static/upload 上传目录）
	r.Static("/static", "./static")

	// metrics endpoint (Prometheus)
	//r.GET("/metrics", gin.WrapH(infraProm.Handler()))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health/ping", health.Ping)

		// 注册路由
		router.NewAIRouter(v1, aiChat, auth)
		router.NewUserRouter(v1, userCtrl, auth)
		router.NewCategoryRouter(v1, categoryCtrl, auth)
		router.NewArticleRouter(v1, articleCtrl, auth)
		router.NewStorageRouter(v1, storageCtrl, auth)
		router.NewCommentRouter(v1, commentCtrl, auth)
		router.NewLikeRouter(v1, likeCtrl, auth)
		router.NewAuthRouter(v1, authCtrl, auth)

	}

	return r
}
