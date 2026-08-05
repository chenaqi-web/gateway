package middleware

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (uint64, bool) {
	userID, ok := c.Get(AuthUserIDContextKey)
	if !ok {
		return 0, false
	}
	return userID.(uint64), true
}
