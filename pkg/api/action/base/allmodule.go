package base

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/base"
)

func (b *baseServer) CreateGroup(g *gin.Context)  {
	request := &apiResource.ModuleRequest{}
	if err := g.ShouldBindJSON(request); err != nil {
		api.ResponseError(g, err)
		return
	}

	response, err := b.AllModuleService.CreateGroup(request)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, response, "")
}

func (b *baseServer) CreateModule(g *gin.Context)  {
	request := &apiResource.ModuleRequest{}
	if err := g.ShouldBindJSON(request); err != nil {
		api.ResponseError(g, err)
		return
	}

	response, err := b.AllModuleService.CreateModule(request)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, response, "")
}

func (b *baseServer) DeleteGroupAndModule(g *gin.Context) {
	request := &apiResource.ModuleRequest{}
	if err := g.ShouldBindJSON(request); err != nil {
		api.ResponseError(g, err)
		return
	}

	response, err := b.AllModuleService.DeleteGroupAndModule(request.UUID)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"delete": response}, "")
}

func (b *baseServer) ListAll(g *gin.Context)  {
	data, err := b.AllModuleService.ListAll()
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, data, "")
}