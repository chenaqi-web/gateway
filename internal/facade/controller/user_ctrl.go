package controller

import (
	"gateway/internal/config"
	"net/http"

	"gateway/internal/application"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (u *UserController) List(c *gin.Context) {
	if middleware.GetRole(c) != "admin" {
		reponse.Fail(c, http.StatusForbidden, "admin access required")
		return
	}
	var query struct {
		Keyword  string `form:"keyword"`
		Page     uint32 `form:"page"`
		PageSize uint32 `form:"page_size"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}
	users, total, err := u.svc.List(c.Request.Context(), query.Keyword, query.Page, query.PageSize)
	if err != nil {
		userRPCError(c, err)
		return
	}
	reponse.Success(c, gin.H{"users": users, "total": total})
}

func (u *UserController) UpdateStatus(c *gin.Context) {
	if middleware.GetRole(c) != "admin" {
		reponse.Fail(c, http.StatusForbidden, "admin access required")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=approved blocked"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	success, err := u.svc.UpdateStatus(c.Request.Context(), userID, req.Status)
	if err != nil {
		userRPCError(c, err)
		return
	}
	reponse.Success(c, gin.H{"success": success})
}

func (u *UserController) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	result, err := u.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (u *UserController) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	result, err := u.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		userRPCError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (u *UserController) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	result, err := u.svc.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		userRPCError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (u *UserController) UpdateAvatar(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Avatar string `json:"avatar" binding:"required,max=500"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	result, err := u.svc.UpdateAvatar(c.Request.Context(), userID, req.Avatar)
	if err != nil {
		userRPCError(c, err)
		return
	}
	reponse.Success(c, result)
}

func userRPCError(c *gin.Context, err error) {
	rpcStatus := status.Convert(err)
	httpStatus := http.StatusInternalServerError
	switch rpcStatus.Code() {
	case codes.InvalidArgument:
		httpStatus = http.StatusBadRequest
	case codes.NotFound:
		httpStatus = http.StatusNotFound
	case codes.AlreadyExists:
		httpStatus = http.StatusConflict
	case codes.PermissionDenied:
		httpStatus = http.StatusForbidden
	}
	reponse.Fail(c, httpStatus, rpcStatus.Message())
}
