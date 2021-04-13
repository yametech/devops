package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/resource"
	"net/http"
)

func (s *Server) ListAppProject(g *gin.Context) {

	data, count,  err := s.AppProjectService.List(1, 10)
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}
	g.JSON(http.StatusOK, gin.H{"data": data, "count": count})
}

func (s *Server) CreateAppProject(g *gin.Context)  {
	var a resource.AppProject
	g.ShouldBindJSON(&a)
	s.AppProjectService.Create(&a)
	g.JSON(http.StatusOK, gin.H{"data": a})
}