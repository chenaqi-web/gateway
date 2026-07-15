package router

import (
	"backend/gateway/internal/facade/controller"
	"backend/gateway/internal/facade/middleware"

	"github.com/gin-gonic/gin"

	"backend/gateway/internal/config"
)

func New(cfg *config.Config, health *controller.HealthController) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 跨域中间件
	r.Use(middleware.Cors())
	// r.Use(middleware.Recovery(), middleware.CORS())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health/ping", health.Ping)
	}

	return r
}
