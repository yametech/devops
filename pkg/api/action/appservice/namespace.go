package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"strconv"
)

func (s *Server) ListNamespaces(g *gin.Context){
	results, err := s.NamespaceService.List()
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, results)
}

func (s *Server) ListByLevel(g *gin.Context)  {
	level, err := strconv.Atoi(g.DefaultQuery("level", "0"))
	if err != nil {
		api.ResponseError(g, errors.New("the level need int type"))
		return
	}

	results, err := s.NamespaceService.ListByLevel(level)
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