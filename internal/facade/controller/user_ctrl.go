package controller

import (
	"gateway/internal/application"
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
)

type UserController struct{ svc *application.UserService }

func NewUserController(svc *application.UserService) *UserController {
	return &UserController{svc: svc}
}

func (u *UserController) Get(c *gin.Context) {
	resp, err := u.svc.Get(c.Request.Context())
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, resp)
}
