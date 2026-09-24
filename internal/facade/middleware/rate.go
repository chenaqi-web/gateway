package middleware

import (
	"github.com/gin-gonic/gin"
)

func PreMinuteLimit() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}
