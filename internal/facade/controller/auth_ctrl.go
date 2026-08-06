package controller

import (
	"context"
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
	rpc          *rpc.Client
	jwtBlackList *cache.JwtBlacklist
	email        *utils.Email
	cfg          config.AuthConfig
	log          *clog.Log
}

func NewAuthController(
	rpcClient *rpc.Client,
	jwtBlackList *cache.JwtBlacklist,
	cfg *config.Config,
	logger *clog.Log,
) *AuthController {
	return &AuthController{
		rpc:          rpcClient,
		jwtBlackList: jwtBlackList,
		email:        utils.NewEmail(cfg, jwtBlackList.Cache),
		cfg:          cfg.Auth,
		log:          logger,
	}
}

func (a *AuthController) SendEmailCode(c *gin.Context) {
	var request dto.SendEmailCodeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), a.rpc.GetRequestTimeout())
	defer cancel()

	err := a.email.SendCode(ctx, request.Email, request.Purpose)
	if err != nil {
		a.log.Error("send email code", zap.Error(err))
		reponse.Fail(c, http.StatusBadGateway, "core-server unavailable")
		return
	}

	reponse.Success(c, nil)
}

func (a *AuthController) EmailLogin(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), a.rpc.GetRequestTimeout())
	defer cancel()
	if err := a.email.VerifyCode(ctx, request.Email, request.Code); err != nil {
		reponse.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}

	response, err := a.rpc.GetAuthClient().EmailLogin(ctx, &authpb.EmailLoginRequest{
		Email: request.Email,
	})
	if err != nil || response == nil || response.GetUser() == nil {
		reponse.Fail(c, http.StatusBadGateway, "invalid core response or error")
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
		if err := a.jwtBlackList.BlacklistToken(c.Request.Context(), authorization[1], a.cfg.AccessExpire); err != nil {
			a.log.Error("auth blacklist access token", zap.Error(err))
			reponse.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	if refreshToken != "" {
		if err := a.jwtBlackList.BlacklistToken(c.Request.Context(), refreshToken, a.cfg.RefreshExpire); err != nil {
			a.log.Error("auth blacklist refresh token", zap.Error(err))
			reponse.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	reponse.Success(c, nil)
}
