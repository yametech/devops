package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
)

func (s *Server) ListNamespaces(g *gin.Context){
	results, err := s.NamespaceService.List()
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, results)
}

func (s *Server) CreateNamespace(g *gin.Context) {
	req := &apiResource.Request{}
	if err := g.ShouldBindJSON(&req); err != nil {
		api.ResponseError(g, err)
		return
	}

	namespace, err := s.NamespaceService.Create(req)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, namespace)
}