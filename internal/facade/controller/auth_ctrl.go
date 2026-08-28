package controller

import (
	"gateway/internal/application"
	"gateway/internal/config"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"gateway/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	cfg *config.Config
	svc *application.AuthService
}

func NewAuthController(svc *application.AuthService, cfg *config.Config) *AuthController {
	return &AuthController{
		cfg: cfg,
		svc: svc,
	}
}

func (a *AuthController) SendEmailCode(c *gin.Context) {
	var req dto.SendEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	if err := a.svc.SendEmailCode(c.Request.Context(), req); err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, nil)
}

func (a *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	if err := a.svc.Register(c.Request.Context(), req); err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, nil)
}

func (a *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)

		return
	}

	result, refreshToken, err := a.svc.Login(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}

	// 在cookie设置refresh_token
	utils.SetRefreshCookie(c.Writer, refreshToken, a.cfg.Auth)
	reponse.Success(c, result)
}

func (a *AuthController) EmailLogin(c *gin.Context) {
	var req dto.EmailLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)

		return
	}
	result, refreshToken, err := a.svc.EmailLogin(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	utils.SetRefreshCookie(c.Writer, refreshToken, a.cfg.Auth)
	reponse.Success(c, result)
}

func (a *AuthController) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)

		return
	}
	if err := a.svc.ForgotPassword(c.Request.Context(), req); err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, nil)
}

func (a *AuthController) Logout(c *gin.Context) {
	refreshToken, _ := utils.RefreshTokenFromCookie(c.Request)
	utils.ClearRefreshCookie(c.Writer, a.cfg.Auth)
	if err := a.svc.Logout(c.Request.Context(), c.GetHeader("Authorization"), refreshToken); err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, nil)
}
