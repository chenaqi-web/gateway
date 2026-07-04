package router

import (
	"github.com/gin-gonic/gin"

	"backend/gateway/internal/config"
)

// New builds the Gin engine and registers all HTTP routes.
func New(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()
	// r.Use(middleware.Recovery(), middleware.CORS())

	v1 := r.Group("/api/v1")
	{
		_ = v1 // register controllers here
	}

	return r
}
