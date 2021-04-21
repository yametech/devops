package allconfigservice

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/resource"
	"net/http"
)

func (s *Server) ListAllConfigProject(g *gin.Context) {
	search := g.Query("search")
	uuid := g.Query("uuid")
	count, name, data, err := s.AllConfigService.List(search, uuid)
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}
	g.JSON(http.StatusOK, gin.H{"count": count, "name": name, "content": data})
}

func (s *Server) CreateAllConfigProject(g *gin.Context) {
	var build resource.AllConfigProject
	uuid, err := s.AllConfigService.Create(&build)
	if err != nil {
		api.RequestParamsError(g, "error", err)
	}
	g.JSON(http.StatusOK, gin.H{"uuid": uuid})
}

func (s *Server) UpdateAllConfigProject(g *gin.Context) {
	uuid := g.Param("uuid")
	name := g.PostForm("name")
	value := g.PostForm("value")
	var build resource.AllConfigProject
	if err := g.ShouldBind(&build); err != nil {
		api.RequestParamsError(g, "error", err)
	}

	data, update, err := s.AllConfigService.Update(name, uuid, value, &build)
	if err != nil {
		api.RequestParamsError(g, "error", err)
	}
	g.JSON(http.StatusOK, gin.H{"data": data, "update": update})
}

func (s *Server) DeleteAllConfigProject(g *gin.Context) {
	uuid := g.Param("uuid")
	name := g.PostForm("name")
	//var build resource.AllConfigProject
	//if err:=g.ShouldBind(&build);err!=nil {
	//	api.RequestParamsError(g,"error",err)
	//}
	err := s.AllConfigService.Delete(uuid, name)
	if err != nil {
		api.RequestParamsError(g, "error", err)
	}
	g.JSON(http.StatusOK, gin.H{"delete": true})
}
