package controller

import (
	"context"
	"log"
	"net/http"
	"strings"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/authpb"
	"gateway/internal/config"
	"gateway/internal/infras/cache"
	"gateway/internal/infras/clog"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"gateway/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthController struct {
	rpc   *rpc.Client
	cache *cache.CacheClient
	cfg   config.AuthConfig
	log   *clog.Log
}

func NewAuthController(
	rpcClient *rpc.Client,
	cacheClient *cache.CacheClient,
	cfg *config.Config,
	logger *clog.Log,
) *AuthController {
	return &AuthController{
		rpc:   rpcClient,
		cache: cacheClient,
		cfg:   cfg.Auth,
		log:   logger,
	}
}

func (a *AuthController) SendEmailCode(c *gin.Context) {
	var request dto.SendEmailCodeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var purpose authpb.EmailCodePurpose
	switch strings.TrimSpace(request.Purpose) {
	case "register":
		purpose = authpb.EmailCodePurpose_EMAIL_CODE_PURPOSE_REGISTER
	case "reset_password":
		purpose = authpb.EmailCodePurpose_EMAIL_CODE_PURPOSE_RESET_PASSWORD
	default:
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), a.rpc.GetRequestTimeout())
	defer cancel()

	response, err := a.rpc.GetAuthClient().SendEmailCode(ctx, &authpb.SendEmailCodeRequest{
		Email:   request.Email,
		Purpose: purpose,
	})
	if err != nil {
		log.Printf("auth send email code: %v", err)
		reponse.Fail(c, http.StatusBadGateway, "core-server unavailable")
		return
	}
	if response == nil || !response.GetSuccess() {
		reponse.Fail(c, http.StatusBadGateway, "invalid core response")
		return
	}

	reponse.Success(c, nil)
}

func (a *AuthController) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), a.rpc.GetRequestTimeout())
	defer cancel()

	// 调用Core登录
	response, err := a.rpc.GetAuthClient().Login(ctx, &authpb.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		reponse.Fail(c, http.StatusBadGateway, err.Error())
		return
	}

	user := response.GetUser()
	// Gateway生成token
	claims := utils.JWTClaims{
		UserID:      user.GetId(),
		Role:        user.GetRole(),
		AuthVersion: user.GetAuthVersion(),
	}

	accessToken, err := utils.CreateAccessToken([]byte(a.cfg.JWTSecret), claims, a.cfg.AccessExpire)
	if err != nil {
		a.log.Error("auth create access token", zap.Error(err))
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}
	refreshToken, err := utils.CreateRefreshToken([]byte(a.cfg.JWTSecret), claims, a.cfg.RefreshExpire)
	if err != nil {
		a.log.Error("auth create refresh token", zap.Error(err))
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}
	// refresh token写入cookie
	utils.SetRefreshCookie(c.Writer, refreshToken, a.cfg)

	reponse.Success(c, dto.LoginResponse{
		AccessToken:     accessToken,
		AccessExpiresIn: a.cfg.AccessExpire,
		User: &dto.AuthUser{
			ID:          user.GetId(),
			Username:    user.GetUsername(),
			Email:       user.GetEmail(),
			Phone:       user.GetPhone(),
			Avatar:      user.GetAvatar(),
			Sex:         user.GetSex(),
			Age:         user.GetAge(),
			Role:        user.GetRole(),
			Status:      user.GetStatus(),
			AuthVersion: user.GetAuthVersion(),
		},
	})
}

func (a *AuthController) Logout(c *gin.Context) {
	refreshToken, _ := utils.RefreshTokenFromCookie(c.Request)
	utils.ClearRefreshCookie(c.Writer, a.cfg)

	// token加入黑名单
	authorization := strings.Fields(c.GetHeader("Authorization"))
	if len(authorization) == 2 && strings.EqualFold(authorization[0], "Bearer") {
		if err := a.cache.BlacklistToken(c.Request.Context(), authorization[1], a.cfg.AccessExpire); err != nil {
			a.log.Error("auth blacklist access token", zap.Error(err))
			reponse.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	if refreshToken != "" {
		if err := a.cache.BlacklistToken(c.Request.Context(), refreshToken, a.cfg.RefreshExpire); err != nil {
			a.log.Error("auth blacklist refresh token", zap.Error(err))
			reponse.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	reponse.Success(c, nil)
}
