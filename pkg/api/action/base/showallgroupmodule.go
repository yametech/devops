package base

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/api"
)

func (b *baseServer) ListGroup(c *gin.Context) {
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		api.ResponseError(c, errors.New("page need int type"))
		return
	}
	pageSize, err := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 64)
	if err != nil {
		api.ResponseError(c, errors.New("pageSize need int type"))
		return
	}
	response, err := b.ListAllGroup(page, pageSize)
	if err != nil {
		api.ResponseError(c, err)
		return
	}
	api.ResponseSuccess(c, response, "")
}

func (b *baseServer) ListModule(c *gin.Context) {
	uuid := c.Query("uuid")
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		api.ResponseError(c, errors.New("page need int type"))
		return
	}
	pageSize, err := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 64)
	if err != nil {
		api.ResponseError(c, errors.New("pageSize need int type"))
		return
	}
	response, err := b.ListAllModule(uuid, page, pageSize)
	if err != nil {
		api.ResponseError(c, err)
		return
	}
	api.ResponseSuccess(c, response, "")
}
