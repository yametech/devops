package base

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/base"
)

func (b *baseServer) AddCollectionModule(g *gin.Context)  {
	request := &apiResource.ModuleRequest{}
	if err := g.ShouldBindJSON(request); err != nil {
		api.ResponseError(g, err)
		return
	}

	response, err := b.CollectionModuleService.AddCollectionModule(request.UUID)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, response, "")
}

func (b *baseServer) ListCollectionModule(g *gin.Context)  {

	response, err := b.CollectionModuleService.ListCollectionModule()
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, response, "")
}

func (b *baseServer) DeleteCollectionModule(g *gin.Context) {
	request := &apiResource.ModuleRequest{}
	if err := g.ShouldBindJSON(request); err != nil {
		api.ResponseError(g, err)
		return
	}

	response, err := b.CollectionModuleService.DeleteCollectionModule(request.UUID)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, response, "")
}
