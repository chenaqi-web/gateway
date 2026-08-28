package controller

import (
	"gateway/internal/application"
	"gateway/internal/config"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	svc *application.UserService
	cfg *config.Config
}

func NewUserController(svc *application.UserService, cfg *config.Config) *UserController {
	return &UserController{
		svc: svc,
		cfg: cfg,
	}
}

func (u *UserController) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	result, err := u.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}

func (u *UserController) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	result, err := u.svc.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}

func (u *UserController) UpdateAvatar(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}

	var req struct {
		Avatar string `json:"avatar" binding:"required,max=500"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	result, err := u.svc.UpdateAvatar(c.Request.Context(), userID, req.Avatar)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}

// =====================================================================================================================
// 用户管理方面

func (u *UserController) UserList(c *gin.Context) {
	if middleware.GetRole(c) != "admin" {
		reponse.Forbidden(c)
		return
	}
	var rep dto.UserListFormRequest
	if err := c.ShouldBindQuery(&rep); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	users, total, err := u.svc.UserList(c.Request.Context(), rep.Keyword, rep.Page, rep.PageSize)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, gin.H{"users": users, "total": total})
}

func (u *UserController) UpdateStatus(c *gin.Context) {
	if middleware.GetRole(c) != "admin" {
		reponse.Forbidden(c)
		return
	}

	var req dto.UpdateUserBlacklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	success, err := u.svc.UpdateBlacklist(c.Request.Context(), req.UserID, req.Blacklisted)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, gin.H{"success": success})
}
