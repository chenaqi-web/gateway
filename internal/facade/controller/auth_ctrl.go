package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"backend/gateway/internal/client/rpc"
	"backend/gateway/internal/client/rpc/core-rpc/authpb"
	"backend/gateway/internal/config"
	"backend/gateway/internal/model/dto"
	"backend/gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authRPC interface {
	SendEmailCode(context.Context, *authpb.SendEmailCodeRequest, ...grpc.CallOption) (*authpb.SendEmailCodeResponse, error)
	Register(context.Context, *authpb.RegisterRequest, ...grpc.CallOption) (*authpb.RegisterResponse, error)
	Login(context.Context, *authpb.LoginRequest, ...grpc.CallOption) (*authpb.LoginResponse, error)
	EmailLogin(context.Context, *authpb.EmailLoginRequest, ...grpc.CallOption) (*authpb.LoginResponse, error)
	RefreshToken(context.Context, *authpb.RefreshTokenRequest, ...grpc.CallOption) (*authpb.RefreshTokenResponse, error)
	Logout(context.Context, *authpb.LogoutRequest, ...grpc.CallOption) (*authpb.LogoutResponse, error)
	ResetPasswordByEmail(context.Context, *authpb.ResetPasswordByEmailRequest, ...grpc.CallOption) (*authpb.ResetPasswordByEmailResponse, error)
}

type AuthController struct {
	client         authRPC
	requestTimeout time.Duration
	cookies        *RefreshCookieManager
}

func NewAuthController(rpcClient *rpc.Client, cfg *config.Config) (*AuthController, error) {
	if rpcClient == nil || cfg == nil {
		return nil, fmt.Errorf("auth controller dependency is nil")
	}
	cookies, err := NewRefreshCookieManager(cfg.Auth)
	if err != nil {
		return nil, err
	}
	return newAuthController(rpcClient.GetAuthClient(), rpcClient.GetRequestTimeout(), cookies)
}

func newAuthController(client authRPC, timeout time.Duration, cookies *RefreshCookieManager) (*AuthController, error) {
	if client == nil || timeout <= 0 || cookies == nil {
		return nil, fmt.Errorf("auth controller dependency is invalid")
	}
	return &AuthController{client: client, requestTimeout: timeout, cookies: cookies}, nil
}

func (a *AuthController) SendEmailCode(c *gin.Context) {
	var request dto.SendEmailCodeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeInvalidAuthRequest(c)
		return
	}
	purpose, ok := emailCodePurpose(request.Purpose)
	if !ok {
		writeInvalidAuthRequest(c)
		return
	}

	ctx, cancel := a.rpcContext(c)
	defer cancel()
	response, err := a.client.SendEmailCode(ctx, &authpb.SendEmailCodeRequest{Email: request.Email, Purpose: purpose})
	if err != nil {
		writeAuthRPCError(c, authOperationGeneric, err)
		return
	}
	if response == nil || !response.GetSuccess() {
		writeAuthRPCError(c, authOperationGeneric, status.Error(codes.Internal, "invalid core response"))
		return
	}
	c.JSON(http.StatusOK, reponse.Success(nil))
}

func (a *AuthController) Register(c *gin.Context) {
	var request dto.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeInvalidAuthRequest(c)
		return
	}

	ctx, cancel := a.rpcContext(c)
	defer cancel()
	response, err := a.client.Register(ctx, &authpb.RegisterRequest{
		Username: request.Username, Email: request.Email, Password: request.Password,
		ConfirmPassword: request.ConfirmPassword, EmailCode: request.EmailCode,
		Phone: request.Phone, Avatar: request.Avatar, Sex: request.Sex, Age: request.Age,
	})
	if err != nil {
		writeAuthRPCError(c, authOperationRegister, err)
		return
	}
	if response == nil || !response.GetSuccess() || response.GetUser() == nil {
		writeAuthRPCError(c, authOperationGeneric, status.Error(codes.Internal, "invalid core response"))
		return
	}
	c.JSON(http.StatusOK, reponse.Success(dto.RegisterResponse{User: authUserSummary(response.GetUser())}))
}

func (a *AuthController) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeInvalidAuthRequest(c)
		return
	}
	a.handleLogin(c, func(ctx context.Context) (*authpb.LoginResponse, error) {
		return a.client.Login(ctx, &authpb.LoginRequest{Username: request.Username, Password: request.Password})
	})
}

func (a *AuthController) EmailLogin(c *gin.Context) {
	var request dto.EmailLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeInvalidAuthRequest(c)
		return
	}
	a.handleLogin(c, func(ctx context.Context) (*authpb.LoginResponse, error) {
		return a.client.EmailLogin(ctx, &authpb.EmailLoginRequest{Email: request.Email, Password: request.Password})
	})
}

func (a *AuthController) Refresh(c *gin.Context) {
	refreshToken, err := a.cookies.Get(c.Request)
	if err != nil {
		a.cookies.Clear(c.Writer)
		writeAuthUnauthenticated(c, authOperationRefresh)
		return
	}

	ctx, cancel := a.rpcContext(c)
	defer cancel()
	response, err := a.client.RefreshToken(ctx, &authpb.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		a.cookies.Clear(c.Writer)
		switch status.Code(err) {
		case codes.Unavailable, codes.DeadlineExceeded:
			writeAuthRPCError(c, authOperationRefresh, err)
		case codes.Unauthenticated, codes.PermissionDenied:
			writeAuthUnauthenticated(c, authOperationRefresh)
		default:
			writeAuthRPCError(c, authOperationRefresh, err)
		}
		return
	}
	tokens := response.GetTokens()
	if tokens == nil || tokens.GetAccessToken() == "" || tokens.GetRefreshToken() == "" || tokens.GetAccessExpiresIn() <= 0 {
		a.cookies.Clear(c.Writer)
		writeAuthRPCError(c, authOperationGeneric, status.Error(codes.Internal, "invalid core response"))
		return
	}
	a.cookies.Set(c.Writer, tokens.GetRefreshToken())
	c.JSON(http.StatusOK, reponse.Success(dto.RefreshResponse{
		AccessToken: tokens.GetAccessToken(), AccessExpiresIn: tokens.GetAccessExpiresIn(),
	}))
}

func (a *AuthController) Logout(c *gin.Context) {
	refreshToken, cookieErr := a.cookies.Get(c.Request)
	a.cookies.Clear(c.Writer)
	if cookieErr != nil {
		c.JSON(http.StatusOK, reponse.Success(nil))
		return
	}

	ctx, cancel := a.rpcContext(c)
	defer cancel()
	_, err := a.client.Logout(ctx, &authpb.LogoutRequest{RefreshToken: refreshToken})
	if err != nil && status.Code(err) != codes.Unauthenticated && status.Code(err) != codes.NotFound {
		writeAuthRPCError(c, authOperationRefresh, err)
		return
	}
	c.JSON(http.StatusOK, reponse.Success(nil))
}

func (a *AuthController) ResetPasswordByEmail(c *gin.Context) {
	var request dto.ResetPasswordByEmailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeInvalidAuthRequest(c)
		return
	}

	ctx, cancel := a.rpcContext(c)
	defer cancel()
	response, err := a.client.ResetPasswordByEmail(ctx, &authpb.ResetPasswordByEmailRequest{
		Email: request.Email, EmailCode: request.EmailCode,
		NewPassword: request.NewPassword, ConfirmPassword: request.ConfirmPassword,
	})
	if err != nil {
		writeAuthRPCError(c, authOperationGeneric, err)
		return
	}
	if response == nil || !response.GetSuccess() {
		writeAuthRPCError(c, authOperationGeneric, status.Error(codes.Internal, "invalid core response"))
		return
	}
	a.cookies.Clear(c.Writer)
	c.JSON(http.StatusOK, reponse.Success(nil))
}

func (a *AuthController) handleLogin(c *gin.Context, call func(context.Context) (*authpb.LoginResponse, error)) {
	ctx, cancel := a.rpcContext(c)
	defer cancel()
	response, err := call(ctx)
	if err != nil {
		writeAuthRPCError(c, authOperationLogin, err)
		return
	}
	tokens := response.GetTokens()
	if response.GetUser() == nil || tokens == nil || tokens.GetAccessToken() == "" ||
		tokens.GetRefreshToken() == "" || tokens.GetAccessExpiresIn() <= 0 {
		writeAuthRPCError(c, authOperationGeneric, status.Error(codes.Internal, "invalid core response"))
		return
	}
	a.cookies.Set(c.Writer, tokens.GetRefreshToken())
	c.JSON(http.StatusOK, reponse.Success(dto.LoginResponse{
		AccessToken: tokens.GetAccessToken(), AccessExpiresIn: tokens.GetAccessExpiresIn(),
		User: authUserSummary(response.GetUser()),
	}))
}

func (a *AuthController) rpcContext(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), a.requestTimeout)
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
		ID: user.GetId(), Username: user.GetUsername(), Email: user.GetEmail(),
		Avatar: user.GetAvatar(), Role: user.GetRole(),
	}
}
