package controller

import (
	"context"
	"gateway/internal/infras/clog"
	"log"
	"net/http"
	"strings"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/authpb"
	"gateway/internal/config"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"gateway/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthController struct {
	rpc *rpc.Client
	cfg config.AuthConfig
	log *clog.Log
}

func NewAuthController(rpcClient *rpc.Client, cfg *config.Config, logger *clog.Log) *AuthController {
	return &AuthController{
		rpc: rpcClient,
		cfg: cfg.Auth,
		log: logger,
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

	response, err := a.rpc.GetAuthClient().Login(ctx, &authpb.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		reponse.Fail(c, http.StatusBadGateway, err.Error())
		return
	}

	user := response.GetUser()
	claims := utils.JWTClaims{
		UserID:      user.GetId(),
		Role:        user.GetRole(),
		AuthVersion: user.GetAuthVersion(),
	}

	accessToken, err := utils.CreateAccessToken(a.cfg.JWTSigningKey(), claims, a.cfg.AccessDuration())
	if err != nil {
		a.log.Error("auth create access token: %v", zap.Error(err))
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}
	refreshToken, err := utils.CreateRefreshToken(a.cfg.JWTSigningKey(), claims, a.cfg.RefreshDuration())
	if err != nil {
		a.log.Error("auth create refresh token: %v", zap.Error(err))
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}
	if err := utils.SetRefreshCookie(c.Writer, refreshToken, a.cfg); err != nil {
		a.log.Error("auth login cookie: %v", zap.Error(err))
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}

	reponse.Success(c, dto.LoginResponse{
		AccessToken:     accessToken,
		AccessExpiresIn: int64(a.cfg.AccessDuration()),
		User: &dto.AuthUserSummary{
			ID:       user.GetId(),
			Username: user.GetUsername(),
			Email:    user.GetEmail(),
			Avatar:   user.GetAvatar(),
			Role:     user.GetRole(),
		},
	})
}

func (a *AuthController) Logout(c *gin.Context) {
	utils.ClearRefreshCookie(c.Writer, a.cfg)
	reponse.Success(c, nil)
}
