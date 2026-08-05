package middleware

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (uint64, bool) {
	userID, ok := c.Get("userID")
	if !ok {
		return 1, false
	}
	return userID.(uint64), true
}
