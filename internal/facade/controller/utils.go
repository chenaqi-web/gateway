package controller

import (
	"gateway/internal/model/reponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

func bindJSON(c *gin.Context, value any) bool {
	if err := c.ShouldBindJSON(value); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return false
	}
	return true
}

func rpcError(c *gin.Context, err error) {
	reponse.Fail(c, http.StatusInternalServerError, err.Error())
}
