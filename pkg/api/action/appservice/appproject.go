package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/apppservice"
)

func (s *Server) ListAppProject(g *gin.Context) {
	search := g.Query("search")

	data, err := s.AppProjectService.List(search)
	if err != nil {
		api.ResponseError(g, err)
		return
	}
	api.ResponseSuccess(g, data)
}

func (s *Server) CreateAppProject(g *gin.Context) {

	request := &apiResource.Request{}
	if err := g.ShouldBindJSON(&request); err != nil {
		api.ResponseError(g, err)
		return
	}

	req, err := s.AppProjectService.Create(request)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"results": req})
}

func (s *Server) UpdateAppProject(g *gin.Context) {
	uuid := g.Param("uuid")
	req := &apiResource.Request{}
	if err := g.ShouldBindJSON(&req); err != nil {
		api.ResponseError(g, err)
		return
	}

	data, update, err := s.AppProjectService.Update(uuid, req)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"results": data, "update": update})
}

func (s *Server) DeleteAppProject(g *gin.Context) {
	uuid := g.Param("uuid")
	result, err := s.AppProjectService.Delete(uuid)
	if err != nil {
		api.ResponseError(g, err)
		return
	}
	api.ResponseSuccess(g, gin.H{"delete": result})
}
