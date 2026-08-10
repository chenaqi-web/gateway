package controller

import (
	"net/http"

	"gateway/internal/application"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserController struct{ svc *application.UserService }

func NewUserController(svc *application.UserService) *UserController {
	return &UserController{svc: svc}
}

func (u *UserController) Get(c *gin.Context) {
	result, err := u.svc.Get(c.Request.Context())
	if err != nil {
		rpcError(c, err)
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
	if !bindJSON(c, &req) {
		return
	}
	result, err := u.svc.UpdateProfile(c.Request.Context(), userID, req)
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
