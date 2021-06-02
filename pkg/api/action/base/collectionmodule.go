package base

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/base"
)

func (b *baseServer) AddCollectionModule(g *gin.Context)  {
	// Get the user
	user := g.Request.Header.Get("x-wrapper-username")

	request := &apiResource.ModuleRequest{}
	if err := g.ShouldBindJSON(request); err != nil {
		api.ResponseError(g, err)
		return
	}

	response, _, err := b.CollectionModuleService.AddCollectionModule(request.UUID, user)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, response, "")
}

func (b *baseServer) ListCollectionModule(g *gin.Context)  {
	// Get the user
	user := g.Request.Header.Get("x-wrapper-username")

	response, err := b.CollectionModuleService.ListCollectionModule(user)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, response, "")
}

func (b *baseServer) DeleteCollectionModule(g *gin.Context) {
	// Get the user
	user := g.Request.Header.Get("x-wrapper-username")

	request := &apiResource.ModuleRequest{}
	if err := g.ShouldBindJSON(request); err != nil {
		api.ResponseError(g, err)
		return
	}

	response, _, err := b.CollectionModuleService.DeleteCollectionModule(request.UUID, user)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, response, "")
}
