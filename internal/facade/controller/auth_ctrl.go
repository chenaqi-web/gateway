package controller

import (
	"errors"
	"gateway/internal/application"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"gateway/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthController struct {
	svc *application.AuthService
}

func NewAuthController(svc *application.AuthService) *AuthController {
	return &AuthController{svc: svc}
}

func (a *AuthController) SendEmailCode(c *gin.Context) {
	var req dto.SendEmailCodeRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := a.svc.SendEmailCode(c.Request.Context(), req); err != nil {
		reponse.Fail(c, http.StatusBadGateway, "core-server unavailable")
		return
	}
	reponse.Success(c, nil)
}

func (a *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := a.svc.Register(c.Request.Context(), req); err != nil {
		if errors.Is(err, application.ErrInvalidOrExpiredVerificationCode) {
			reponse.Fail(c, http.StatusUnauthorized, err.Error())
		} else {
			reponse.Fail(c, http.StatusBadRequest, err.Error())
		}
		return
	}
	reponse.Success(c, nil)
}

func (a *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !bindJSON(c, &req) {
		return
	}
	result, refreshToken, err := a.svc.Login(c.Request.Context(), req)
	if err != nil {
		authRPCError(c, err)
		return
	}
	utils.SetRefreshCookie(c.Writer, refreshToken, a.svc.RefreshCookieConfig())
	reponse.Success(c, result)
}

func (a *AuthController) EmailLogin(c *gin.Context) {
	var req dto.EmailLoginRequest
	if !bindJSON(c, &req) {
		return
	}
	result, refreshToken, err := a.svc.EmailLogin(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, application.ErrInvalidOrExpiredVerificationCode) {
			reponse.Fail(c, http.StatusUnauthorized, err.Error())
		} else {
			reponse.Fail(c, http.StatusBadGateway, err.Error())
		}
		return
	}
	utils.SetRefreshCookie(c.Writer, refreshToken, a.svc.RefreshCookieConfig())
	reponse.Success(c, result)
}

func (a *AuthController) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := a.svc.ForgotPassword(c.Request.Context(), req); err != nil {
		if errors.Is(err, application.ErrPasswordsDoNotMatch) {
			reponse.Fail(c, http.StatusBadRequest, err.Error())
		} else if errors.Is(err, application.ErrInvalidOrExpiredVerificationCode) {
			reponse.Fail(c, http.StatusUnauthorized, err.Error())
		} else {
			authRPCError(c, err)
		}
		return
	}
	reponse.Success(c, nil)
}

func (a *AuthController) Logout(c *gin.Context) {
	refreshToken, _ := utils.RefreshTokenFromCookie(c.Request)
	utils.ClearRefreshCookie(c.Writer, a.svc.RefreshCookieConfig())
	if err := a.svc.Logout(c.Request.Context(), c.GetHeader("Authorization"), refreshToken); err != nil {
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}
	reponse.Success(c, nil)
}

func authRPCError(c *gin.Context, err error) {
	switch status.Code(err) {
	case codes.InvalidArgument:
		reponse.Fail(c, http.StatusBadRequest, status.Convert(err).Message())
	case codes.Unauthenticated:
		reponse.Fail(c, http.StatusUnauthorized, status.Convert(err).Message())
	case codes.PermissionDenied:
		reponse.Fail(c, http.StatusForbidden, status.Convert(err).Message())
	case codes.Unavailable, codes.DeadlineExceeded:
		reponse.Fail(c, http.StatusBadGateway, "core-server unavailable")
	default:
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
	}
}
