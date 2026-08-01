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

	purpose, ok := emailCodePurpose(request.Purpose)
	if !ok {
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
		writeAuthRPCError(c, err)
		return
	}
	if response == nil || !response.GetSuccess() {
		writeAuthRPCError(c, status.Error(codes.Internal, "invalid core response"))
		return
	}

	reponse.Success(c, nil)
}

func (a *AuthController) Register(c *gin.Context) {
	var request dto.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), a.rpc.GetRequestTimeout())
	defer cancel()

	response, err := a.rpc.GetAuthClient().Register(ctx, &authpb.RegisterRequest{
		Username:        request.Username,
		Email:           request.Email,
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
		EmailCode:       request.EmailCode,
		Phone:           request.Phone,
		Avatar:          request.Avatar,
		Sex:             request.Sex,
		Age:             request.Age,
	})
	if err != nil {
		writeAuthRPCError(c, err)
		return
	}
	if response == nil || !response.GetSuccess() || response.GetUser() == nil {
		writeAuthRPCError(c, status.Error(codes.Internal, "invalid core response"))
		return
	}

	reponse.Success(c, dto.RegisterResponse{User: authUserSummary(response.GetUser())})
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
		writeAuthRPCError(c, err)
		return
	}

	tokens := response.GetTokens()
	if response.GetUser() == nil || tokens == nil || tokens.GetAccessToken() == "" ||
		tokens.GetRefreshToken() == "" || tokens.GetAccessExpiresIn() <= 0 {
		writeAuthRPCError(c, status.Error(codes.Internal, "invalid core response"))
		return
	}
	if err := utils.SetRefreshCookie(c.Writer, tokens.GetRefreshToken(), a.cfg); err != nil {
		log.Printf("auth login cookie: %v", err)
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}

	reponse.Success(c, dto.LoginResponse{
		AccessToken:     tokens.GetAccessToken(),
		AccessExpiresIn: tokens.GetAccessExpiresIn(),
		User:            authUserSummary(response.GetUser()),
	})
}

func (a *AuthController) EmailLogin(c *gin.Context) {
	var request dto.EmailLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), a.rpc.GetRequestTimeout())
	defer cancel()

	response, err := a.rpc.GetAuthClient().EmailLogin(ctx, &authpb.EmailLoginRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		writeAuthRPCError(c, err)
		return
	}

	tokens := response.GetTokens()
	if response.GetUser() == nil || tokens == nil || tokens.GetAccessToken() == "" ||
		tokens.GetRefreshToken() == "" || tokens.GetAccessExpiresIn() <= 0 {
		writeAuthRPCError(c, status.Error(codes.Internal, "invalid core response"))
		return
	}
	if err := utils.SetRefreshCookie(c.Writer, tokens.GetRefreshToken(), a.cfg); err != nil {
		log.Printf("auth email login cookie: %v", err)
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}

	reponse.Success(c, dto.LoginResponse{
		AccessToken:     tokens.GetAccessToken(),
		AccessExpiresIn: tokens.GetAccessExpiresIn(),
		User:            authUserSummary(response.GetUser()),
	})
}

func (a *AuthController) Refresh(c *gin.Context) {
	refreshToken, err := utils.RefreshTokenFromCookie(c.Request)
	if err != nil {
		utils.ClearRefreshCookie(c.Writer, a.cfg)
		reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), a.rpc.GetRequestTimeout())
	defer cancel()

	response, err := a.rpc.GetAuthClient().RefreshToken(ctx, &authpb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		utils.ClearRefreshCookie(c.Writer, a.cfg)
		if status.Code(err) == codes.Unauthenticated || status.Code(err) == codes.PermissionDenied {
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		writeAuthRPCError(c, err)
		return
	}

	tokens := response.GetTokens()
	if tokens == nil || tokens.GetAccessToken() == "" || tokens.GetRefreshToken() == "" ||
		tokens.GetAccessExpiresIn() <= 0 {
		utils.ClearRefreshCookie(c.Writer, a.cfg)
		writeAuthRPCError(c, status.Error(codes.Internal, "invalid core response"))
		return
	}
	if err := utils.SetRefreshCookie(c.Writer, tokens.GetRefreshToken(), a.cfg); err != nil {
		log.Printf("auth refresh cookie: %v", err)
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}

	reponse.Success(c, dto.RefreshResponse{
		AccessToken:     tokens.GetAccessToken(),
		AccessExpiresIn: tokens.GetAccessExpiresIn(),
	})
}

func (a *AuthController) Logout(c *gin.Context) {
	refreshToken, cookieErr := utils.RefreshTokenFromCookie(c.Request)
	utils.ClearRefreshCookie(c.Writer, a.cfg)
	if cookieErr != nil {
		reponse.Success(c, nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), a.rpc.GetRequestTimeout())
	defer cancel()

	_, err := a.rpc.GetAuthClient().Logout(ctx, &authpb.LogoutRequest{RefreshToken: refreshToken})
	if err != nil && status.Code(err) != codes.Unauthenticated && status.Code(err) != codes.NotFound {
		writeAuthRPCError(c, err)
		return
	}

	reponse.Success(c, nil)
}

func (a *AuthController) ResetPasswordByEmail(c *gin.Context) {
	var request dto.ResetPasswordByEmailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), a.rpc.GetRequestTimeout())
	defer cancel()

	response, err := a.rpc.GetAuthClient().ResetPasswordByEmail(ctx, &authpb.ResetPasswordByEmailRequest{
		Email:           request.Email,
		EmailCode:       request.EmailCode,
		NewPassword:     request.NewPassword,
		ConfirmPassword: request.ConfirmPassword,
	})
	if err != nil {
		writeAuthRPCError(c, err)
		return
	}
	if response == nil || !response.GetSuccess() {
		writeAuthRPCError(c, status.Error(codes.Internal, "invalid core response"))
		return
	}

	utils.ClearRefreshCookie(c.Writer, a.cfg)
	reponse.Success(c, nil)
}

func writeAuthRPCError(c *gin.Context, err error) {
	log.Printf("auth rpc: %v", err)

	switch status.Code(err) {
	case codes.InvalidArgument, codes.FailedPrecondition:
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
	case codes.Unauthenticated:
		reponse.Fail(c, http.StatusUnauthorized, "authentication failed")
	case codes.PermissionDenied:
		reponse.Fail(c, http.StatusForbidden, "permission denied")
	case codes.AlreadyExists:
		reponse.Fail(c, http.StatusConflict, "resource already exists")
	case codes.ResourceExhausted:
		reponse.Fail(c, http.StatusTooManyRequests, "too many requests")
	default:
		reponse.Fail(c, http.StatusBadGateway, "core-server unavailable")
	}
}

func emailCodePurpose(value string) (authpb.EmailCodePurpose, bool) {
	switch strings.TrimSpace(value) {
	case "register":
		return authpb.EmailCodePurpose_EMAIL_CODE_PURPOSE_REGISTER, true
	case "reset_password":
		return authpb.EmailCodePurpose_EMAIL_CODE_PURPOSE_RESET_PASSWORD, true
	default:
		return authpb.EmailCodePurpose_EMAIL_CODE_PURPOSE_UNSPECIFIED, false
	}
}

func authUserSummary(user *authpb.UserInfo) *dto.AuthUserSummary {
	if user == nil {
		return nil
	}

	return &dto.AuthUserSummary{
		ID:       user.GetId(),
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
		Avatar:   user.GetAvatar(),
		Role:     user.GetRole(),
	}
}
