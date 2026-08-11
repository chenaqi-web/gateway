package reponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse 统一响应：code 与 HTTP 状态码一致，无业务自定义错误码。
type APIResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// Success 成功响应（HTTP 200）。
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, APIResponse{Code: http.StatusOK, Msg: "success", Data: data})
}

// Fail 失败响应，code 与 HTTP 状态码一致。
func Fail(c *gin.Context, status int, msg string) {
	c.JSON(status, APIResponse{Code: status, Msg: msg})
}

// Unauthorized 未授权或者无权限
func Unauthorized(c *gin.Context) {
	Fail(c, http.StatusUnauthorized, "unauthorized")
}
