package controller

import (
	"context"
	"log"
	"net/http"
	"strings"

	"backend/gateway/internal/client/rpc"
	"backend/gateway/internal/client/rpc/core-rpc/authpb"
	"backend/gateway/internal/config"
	"backend/gateway/internal/model/dto"
	"backend/gateway/internal/model/reponse"
	"backend/gateway/internal/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthController struct {
	rpc *rpc.Client
	cfg config.AuthConfig
}

func NewAuthController(rpcClient *rpc.Client, cfg *config.Config) *AuthController {
	return &AuthController{
		rpc: rpcClient,
		cfg: cfg.Auth,
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
		log.Printf("auth login: %v", err)
		if status.Code(err) == codes.Unauthenticated {
			reponse.Fail(c, http.StatusUnauthorized, "authentication failed")
			return
		}
		reponse.Fail(c, http.StatusBadGateway, "core-server unavailable")
		return
	}
	if response == nil || response.GetUser() == nil || response.GetUser().GetId() == 0 || response.GetUser().GetRole() == "" {
		reponse.Fail(c, http.StatusBadGateway, "invalid core response")
		return
	}

	signingKey, err := a.cfg.JWTSigningKey()
	if err != nil {
		log.Printf("auth login config: %v", err)
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}
	accessExpiresIn, err := a.cfg.AccessDuration()
	if err != nil {
		log.Printf("auth access duration: %v", err)
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}
	refreshExpiresIn, err := a.cfg.RefreshDuration()
	if err != nil {
		log.Printf("auth refresh duration: %v", err)
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}

	user := response.GetUser()
	claims := utils.JWTClaims{
		UserID:      user.GetId(),
		Role:        user.GetRole(),
		AuthVersion: user.GetAuthVersion(),
	}
	accessToken, err := utils.CreateAccessToken(signingKey, claims, accessExpiresIn)
	if err != nil {
		log.Printf("auth create access token: %v", err)
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}
	refreshToken, err := utils.CreateRefreshToken(signingKey, claims, refreshExpiresIn)
	if err != nil {
		log.Printf("auth create refresh token: %v", err)
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}
	if err := utils.SetRefreshCookie(c.Writer, refreshToken, a.cfg); err != nil {
		log.Printf("auth login cookie: %v", err)
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}

	reponse.Success(c, dto.LoginResponse{
		AccessToken:     accessToken,
		AccessExpiresIn: int64(accessExpiresIn.Seconds()),
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
