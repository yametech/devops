package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/resource"
	"net/http"
)

func (s *Server) ListAppProject(g *gin.Context) {
	search := g.Query("search")

	data, count, err := s.AppProjectService.List(search)
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}
	g.JSON(http.StatusOK, gin.H{"data": data, "count": count})
}

func (s *Server) CreateAppProject(g *gin.Context) {
	var app resource.AppProject
	if err := g.ShouldBindJSON(&app); err != nil {
		api.RequestParamsError(g, "error", err)
	}
	if err := s.AppProjectService.Create(&app); err != nil {
		api.RequestParamsError(g, "error", err)
	}
	g.JSON(http.StatusOK, gin.H{"data": app})
}

func (s *Server) RetrieveAppProject(g *gin.Context) {
	uuid := g.Param("uuid")
	data, err := s.AppProjectService.GetByUUID(uuid)
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}
	g.JSON(http.StatusOK, data)
}

func (s *Server) UpdateAppProject(g *gin.Context) {
	uuid := g.Param("uuid")
	var app resource.AppProject
	if err := g.ShouldBindJSON(&app); err != nil {
		api.RequestParamsError(g, "error", err)
	}
	data, update, err := s.AppProjectService.Update(uuid, &app)
	if err != nil {
		api.RequestParamsError(g, "error", err)
	}
	g.JSON(http.StatusOK, gin.H{"data": data, "update": update})
}

func (s *Server) DeleteAppProject(g *gin.Context) {
	uuid := g.Param("uuid")
	if err := s.AppProjectService.Delete(uuid); err != nil {
		api.RequestParamsError(g, "error", err)
	}
	g.JSON(http.StatusOK, gin.H{"delete": true})
}
