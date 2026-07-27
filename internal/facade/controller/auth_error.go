package controller

import (
	"net/http"
	"strings"

	"backend/gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	authCodeInvalidInput       = 10001
	authCodeInvalidCredentials = 10002
	authCodeUsernameExists     = 10003
	authCodeEmailExists        = 10004
	authCodeEmailCodeInvalid   = 10005
	authCodeAccessInvalid      = 10006
	authCodeActiveSession      = 10007
	authCodeMailRateLimited    = 10010
	authCodeUserDisabled       = 10012
	authCodeRefreshInvalid     = 10013
)

type authOperation uint8

const (
	authOperationGeneric authOperation = iota
	authOperationRegister
	authOperationLogin
	authOperationRefresh
)

func writeInvalidAuthRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, reponse.Error(authCodeInvalidInput, "参数错误"))
}

func writeAuthRPCError(c *gin.Context, operation authOperation, err error) {
	grpcStatus := status.Convert(err)
	switch grpcStatus.Code() {
	case codes.InvalidArgument:
		c.JSON(http.StatusBadRequest, reponse.Error(authCodeInvalidInput, "参数错误"))
	case codes.FailedPrecondition:
		c.JSON(http.StatusBadRequest, reponse.Error(authCodeEmailCodeInvalid, "邮箱验证码错误或过期"))
	case codes.AlreadyExists:
		writeAuthConflict(c, operation, grpcStatus.Message())
	case codes.ResourceExhausted:
		c.JSON(http.StatusTooManyRequests, reponse.Error(authCodeMailRateLimited, "邮件发送频繁"))
	case codes.Unauthenticated:
		writeAuthUnauthenticated(c, operation)
	case codes.PermissionDenied:
		c.JSON(http.StatusForbidden, reponse.Error(authCodeUserDisabled, "用户已禁用或权限不足"))
	case codes.Unavailable, codes.DeadlineExceeded:
		c.JSON(http.StatusServiceUnavailable, reponse.Error(http.StatusServiceUnavailable, "鉴权服务暂不可用"))
	default:
		c.JSON(http.StatusInternalServerError, reponse.Error(http.StatusInternalServerError, "内部错误"))
	}
}

func writeAuthConflict(c *gin.Context, operation authOperation, message string) {
	if operation == authOperationLogin {
		c.JSON(http.StatusConflict, reponse.Error(authCodeActiveSession, "该账号已有有效登录会话"))
		return
	}
	message = strings.ToLower(message)
	if strings.Contains(message, "email") {
		c.JSON(http.StatusConflict, reponse.Error(authCodeEmailExists, "邮箱已存在"))
		return
	}
	c.JSON(http.StatusConflict, reponse.Error(authCodeUsernameExists, "用户名已存在"))
}

func writeAuthUnauthenticated(c *gin.Context, operation authOperation) {
	switch operation {
	case authOperationLogin:
		c.JSON(http.StatusUnauthorized, reponse.Error(authCodeInvalidCredentials, "用户名、邮箱或密码错误"))
	case authOperationRefresh:
		c.JSON(http.StatusUnauthorized, reponse.Error(authCodeRefreshInvalid, "Refresh Token 无效或过期"))
	default:
		c.JSON(http.StatusUnauthorized, reponse.Error(authCodeAccessInvalid, "Access Token 无效或过期"))
	}
}
