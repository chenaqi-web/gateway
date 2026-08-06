package controller

import (
	"context"
	"gateway/internal/application"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"gateway/internal/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5e9)
	defer cancel()
	if err := a.svc.SendEmailCode(ctx, req); err != nil {
		log.Printf("send email code: %v", err)
		reponse.Fail(c, http.StatusBadGateway, "core-server unavailable")
		return
	}
	reponse.Success(c, nil)
}

func (a *AuthController) EmailLogin(c *gin.Context) {
	var req dto.LoginRequest
	if !bindJSON(c, &req) {
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5e9)
	defer cancel()
	result, refreshToken, err := a.svc.EmailLogin(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid or expired") {
			reponse.Fail(c, http.StatusUnauthorized, err.Error())
		} else {
			reponse.Fail(c, http.StatusBadGateway, err.Error())
		}
		return
	}
	utils.SetRefreshCookie(c.Writer, refreshToken, a.svc.Config())
	reponse.Success(c, result)
}

func (a *AuthController) Logout(c *gin.Context) {
	refreshToken, _ := utils.RefreshTokenFromCookie(c.Request)
	utils.ClearRefreshCookie(c.Writer, a.svc.Config())
	accessToken := ""
	authorization := strings.Fields(c.GetHeader("Authorization"))
	if len(authorization) == 2 && strings.EqualFold(authorization[0], "Bearer") {
		accessToken = authorization[1]
	}
	if err := a.svc.Logout(c.Request.Context(), accessToken, refreshToken); err != nil {
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return
	}
	reponse.Success(c, nil)
}
