package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	AuthUserIDContextKey    = "user_id"
	AuthSessionIDContextKey = "session_id"
	AuthRoleContextKey      = "role"
)

type Auth struct {
	requestTimeout time.Duration
	// redis client
}

func NewAuth(requestTimeout time.Duration) *Auth {
	return &Auth{}
}

func (m *Auth) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1.请求来（先提取toekn）
		// 2.先判断accesstoken是不是过期
		// a.过期 ： 用refreshtoken来刷新accesstoken
		// b.没过期就没问题
		// 3. refreshtoken 那就踢出去

	}
}
