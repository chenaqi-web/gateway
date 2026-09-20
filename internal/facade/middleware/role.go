package middleware

import (
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
)

func Role() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只有admin才有权限操作
		if GetRole(c) != "admin" {
			reponse.Forbidden(c)
			return
		}

		c.Next()
	}
}
