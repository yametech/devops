package base

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
)

func (b *baseServer) CreateModuleEntry(c *gin.Context) {
	user := c.Request.Header.Get("x-wrapper-username")
	uuid := c.Query("uuid")
	response, err := b.CreateEntry(user, uuid)
	if err != nil {
		api.ResponseError(c, err)
		return
	}
	api.ResponseSuccess(c, response, "")
}

func (b *baseServer) DeleteModuleEntry(c *gin.Context) {
	user := c.Request.Header.Get("x-wrapper-username")
	uuid := c.Query("uuid")
	response, err := b.DeleteEntry(user, uuid)
	if err != nil {
		api.ResponseError(c, err)
		return
	}
	api.ResponseSuccess(c, response, "")
}

func (b *baseServer) QueryModuleEntry(c *gin.Context) {
	user := c.Request.Header.Get("x-wrapper-username")
	response, err := b.QueryEntry(user)
	if err != nil {
		api.ResponseError(c, err)
		return
	}
	api.ResponseSuccess(c, response, "")
}
