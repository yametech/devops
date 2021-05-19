package base

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/api"
	"strconv"
)

func (b *baseServer) ListGroup(g *gin.Context) {
	page, err := strconv.ParseInt(g.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		api.ResponseError(g, errors.New("page need int type"))
		return
	}
	pageSize, err := strconv.ParseInt(g.DefaultQuery("page_size", "10"), 10, 64)
	if err != nil {
		api.ResponseError(g, errors.New("pageSize need int type"))
		return
	}
	response, err := b.ListAllGroup(page, pageSize)
	if err != nil {
		api.ResponseError(g, err)
		return
	}
	api.ResponseSuccess(g, response, "")
}

func (b *baseServer) ListModule(g *gin.Context) {
	uuid := g.Query("uuid")
	page, err := strconv.ParseInt(g.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		api.ResponseError(g, errors.New("page need int type"))
		return
	}
	pageSize, err := strconv.ParseInt(g.DefaultQuery("page_size", "10"), 10, 64)
	if err != nil {
		api.ResponseError(g, errors.New("pageSize need int type"))
		return
	}
	response, err := b.ListAllModule(uuid, page, pageSize)
	if err != nil {
		api.ResponseError(g, err)
		return
	}
	api.ResponseSuccess(g, response, "")
}
