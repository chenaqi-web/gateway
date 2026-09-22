package controller

import (
	"fmt"
	"gateway/internal/application"
	"gateway/internal/config"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"gateway/internal/utils/vaildate"
	"strconv"

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
		reponse.InternalServerError(c, err.Error())
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
	err := u.svc.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, nil)
}

func (u *UserController) UpdateAvatar(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}

	var req dto.UserAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	req.UserID = userID
	result, err := u.svc.UpdateAvatar(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (u *UserController) UserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	res, err := u.svc.UserList(c.Request.Context(), &dto.UserListRequest{
		Page:     vaildate.Page(page),
		PageSize: vaildate.PageSize(pageSize),
	})
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, res)
}

func (u *UserController) UpdateStatus(c *gin.Context) {
	var req dto.UpdateUserBlacklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	err := u.svc.UpdateBlacklist(c.Request.Context(), req.UserID, req.Blacklisted)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, nil)
}

func (u *UserController) SearchUser(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	result, err := u.svc.SearchUser(c.Request.Context(), &dto.UserSearchRequest{
		Keyword:  keyword,
		Page:     vaildate.Page(page),
		PageSize: vaildate.PageSize(pageSize),
	})
	fmt.Println("ha", result)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}
