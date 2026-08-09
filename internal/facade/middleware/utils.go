package middleware

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (uint64, bool) {
	value, ok := c.Get(AuthUserIDContextKey)
	if !ok {
		return 0, false
	}
	userID, ok := value.(uint64)
	return userID, ok && userID > 0
}

func GetRole(c *gin.Context) string {
	value, ok := c.Get(AuthRoleContextKey)
	if !ok {
		return ""
	}
	role, ok := value.(string)
	if !ok {
		return ""
	}
	return role
}
