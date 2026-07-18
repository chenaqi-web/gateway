package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
// 允许所有来源的跨域请求，支持 GET、POST、PUT、DELETE、OPTIONS 方法
// 允许携带 Authorization 等自定义请求头
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源访问
		c.Header("Access-Control-Allow-Origin", "*")
		// 允许的请求方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 允许的请求头
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-User-Id")

		// 如果是预检请求（OPTIONS），直接返回，不继续处理
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
