package facade

import (
	"backend/gateway/internal/facade/controller"
	"backend/gateway/internal/facade/middleware"
	"backend/gateway/internal/facade/router"

	"github.com/gin-gonic/gin"

	"backend/gateway/internal/config"
)

func New(cfg *config.Config,
	health *controller.HealthController,
	authCtrl *controller.AuthController,
	aiChat *controller.AiChatController,
	userCtrl *controller.UserController,
	categoryCtrl *controller.CategoryController,
	articleCtrl *controller.ArticleController,
	storageCtrl *controller.StorageController,
) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 跨域中间件
	r.Use(middleware.Cors())

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
		router.NewAuthRouter(v1, authCtrl)

		// 注册路由
		router.NewAIRouter(v1, aiChat)
		router.NewUserRouter(v1, userCtrl)
		router.NewCategoryRouter(v1, categoryCtrl)
		router.NewArticleRouter(v1, articleCtrl)
		router.NewStorageRouter(v1, storageCtrl)

	}

	return r
}
