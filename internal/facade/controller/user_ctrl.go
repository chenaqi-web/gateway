package controller

import (
	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/userpb"
	"gateway/internal/model/reponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	rpc *rpc.Client
}

func NewUserController(rpcClient *rpc.Client) *UserController {
	return &UserController{rpc: rpcClient}
}

func (u *UserController) Get(c *gin.Context) {
	c.Get("x-uesr-id")
	resp, err := u.rpc.UserClient.Login(c, &userpb.LoginReq{})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, resp)
}
