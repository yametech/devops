package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
)

func (s *Server) GetNamespaceConfig(g *gin.Context){
	uuid := g.Param("uuid")
	results, err := s.NamespaceConfigService.GetByFilter(uuid)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, results)
}

func (s *Server) UpdateNamespaceConfig(g *gin.Context) {
	data := &apiResource.NameSpaceRequest{}
	if err := g.ShouldBindJSON(&data); err != nil {
		api.ResponseError(g, err)
		return
	}

	result, update, err := s.NamespaceConfigService.Update(data)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"results": result, "update": update})
}