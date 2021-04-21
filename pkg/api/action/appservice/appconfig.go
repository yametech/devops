package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"github.com/yametech/devops/pkg/resource/appproject"
)

func (s *Server) GetAppConfig(g *gin.Context) {
	uuid := g.Param("uuid")
	config, err := s.AppConfigService.GetByFilter(uuid)
	if err != nil {
		api.ResponseSuccess(g, &appproject.AppConfig{})
		return
	}

	api.ResponseSuccess(g, config)
}

func (s *Server) UpdateAppConfig(g *gin.Context) {
	data := &apiResource.AppConfigRequest{}
	if err := g.ShouldBindJSON(&data); err != nil {
		api.ResponseError(g, err)
		return
	}

	result, update, err := s.AppConfigService.Update(data)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"results": result, "update": update})
}
