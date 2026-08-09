package controller

import (
	"gateway/internal/application"
	"gateway/internal/config"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"gateway/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	svc *application.AuthService
	cfg *config.Config
}

func NewAuthController(svc *application.AuthService, cfg *config.Config) *AuthController {
	return &AuthController{
		cfg: cfg,
		svc: svc,
	}
}

func (a *AuthController) SendEmailCode(c *gin.Context) {
	var req dto.SendEmailCodeRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := a.svc.SendEmailCode(c.Request.Context(), req); err != nil {
		reponse.Fail(c, http.StatusInternalServerError, "core-server unavailable")
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
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
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
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 在cookie设置refresh_token
	utils.SetRefreshCookie(c.Writer, refreshToken, a.cfg.Auth)
	reponse.Success(c, result)
}

func (a *AuthController) EmailLogin(c *gin.Context) {
	var req dto.EmailLoginRequest
	if !bindJSON(c, &req) {
		return
	}
	result, refreshToken, err := a.svc.EmailLogin(c.Request.Context(), req)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SetRefreshCookie(c.Writer, refreshToken, a.cfg.Auth)
	reponse.Success(c, result)
}

func (a *AuthController) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := a.svc.ForgotPassword(c.Request.Context(), req); err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, nil)
}

func (a *AuthController) Logout(c *gin.Context) {
	refreshToken, _ := utils.RefreshTokenFromCookie(c.Request)
	utils.ClearRefreshCookie(c.Writer, a.cfg.Auth)
	if err := a.svc.Logout(c.Request.Context(), c.GetHeader("Authorization"), refreshToken); err != nil {
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}
	reponse.Success(c, nil)
}
