package base

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/base"
	"strconv"
)

func (b *baseServer) CreateGlobalModule(g *gin.Context) {
	request := &apiResource.ModuleRequest{}
	if err := g.ShouldBindJSON(request); err != nil {
		api.ResponseError(g, err)
		return
	}

	response, err := b.ChildModuleService.CreateChildModule(request)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, response, "")
}

func (b *baseServer) ListGlobalModule(g *gin.Context) {
	search := g.Query("search")
	parent := g.Query("parent")
	page, err := strconv.ParseInt(g.DefaultQuery("page", "0"), 10, 64)
	if err != nil {
		api.ResponseError(g, errors.New("page need int type"))
		return
	}
	pageSize, err := strconv.ParseInt(g.DefaultQuery("page_size", "0"), 10, 64)
	if err != nil {
		api.ResponseError(g, errors.New("pageSize need int type"))
		return
	}

	response, count, err := b.ChildModuleService.ListChildModule(parent, search, page, pageSize)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"results": response, "count": count}, "")
}

func (b *baseServer) DeleteGlobalModule(g *gin.Context) {
	uuid := g.Param("uuid")
	response, err := b.ChildModuleService.DeleteChildModule(uuid)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, response, "")
}
