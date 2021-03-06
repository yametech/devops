package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/appservice"
)

func (s *Server) GetResourcePoolConfig(g *gin.Context) {
	uuid := g.Param("uuid")
	results, err := s.ResourcePoolConfigService.GetResourcePoolConfig(uuid)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, results, "")
}

func (s *Server) UpdateResourcePoolConfig(g *gin.Context) {
	data := &apiResource.NamespaceRequest{}
	if err := g.ShouldBindJSON(data); err != nil {
		api.ResponseError(g, err)
		return
	}

	result, update, err := s.ResourcePoolConfigService.Update(data)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"results": result, "update": update}, "")
}

func (s *Server) GetNamespaceResourceRemain(g *gin.Context) {
	uuid := g.Param("uuid")
	cpu, memory, err := s.ResourcePoolConfigService.GetNamespaceResourceRemain(uuid)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"cpu_remain": cpu, "memory_remain": memory}, "")
}

func (s *Server) GetNamespaceResource(g *gin.Context){
	uuid := g.Param("uuid")
	cpu, memory, moneyMonth, moneyYear, err := s.ResourcePoolConfigService.GetNamespaceResource(uuid)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"cpu": cpu, "memory": memory, "moneyMonth": moneyMonth, "moneyYear": moneyYear}, "")
}