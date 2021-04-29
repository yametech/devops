package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/apppservice"
)

func (s *Server) GetResourcePoolConfig(g *gin.Context) {
	uuid := g.Param("uuid")
	results, err := s.ResourcePoolConfigService.GetResourcePoolConfig(uuid)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, results)
}

func (s *Server) UpdateResourcePoolConfig(g *gin.Context) {
	data := &apiResource.NamespaceRequest{}
	if err := g.ShouldBindJSON(&data); err != nil {
		api.ResponseError(g, err)
		return
	}

	result, update, err := s.ResourcePoolConfigService.Update(data)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"results": result, "update": update})
}
