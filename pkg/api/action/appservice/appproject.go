package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appproject"
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

	request := &apiResource.AppProjectRequest{}
	if err := g.ShouldBindJSON(&request); err != nil {
		api.ResponseError(g, err)
		return
	}

	req := &appproject.AppProject{
		Metadata: core.Metadata{
			Name: request.Name,
		},
		Spec: appproject.AppSpec{
			AppType: request.AppType,
			ParentApp: request.ParentApp,
			Desc: request.Desc,
			Owner: request.Owner,
		},
	}

	if err := s.AppProjectService.Create(req); err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"results": req})
}

func (s *Server) UpdateAppProject(g *gin.Context) {
	uuid := g.Param("uuid")
	var req apiResource.AppProjectRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		api.ResponseError(g, err)
		return
	}

	app := &appproject.AppProject{
		Spec: appproject.AppSpec{
			Owner: req.Owner,
			Desc: req.Desc,
		},
	}

	data, update, err := s.AppProjectService.Update(uuid, app)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"results": data, "update": update})
}

func (s *Server) DeleteAppProject(g *gin.Context) {
	uuid := g.Param("uuid")
	result, err := s.AppProjectService.Delete(uuid);
	if err != nil {
		api.ResponseError(g, err)
		return
	}
	api.ResponseSuccess(g, gin.H{"delete": result})
}
