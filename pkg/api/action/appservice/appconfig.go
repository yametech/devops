package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"github.com/yametech/devops/pkg/resource/appproject"
)

func (s *Server) GetAppConfig(g *gin.Context) {
	var data apiResource.AppConfigRequest
	if err := g.ShouldBindJSON(&data); err != nil {
		api.ResponseError(g, err)
		return
	}

	config := &appproject.AppConfig{
		Spec: appproject.AppConfigSpec{
			App: data.App,
		},
	}

	if err := s.AppConfigService.GetByFilter(config); err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"results": config})
}

func (s *Server) UpdateAppConfig(g *gin.Context) {
	var data apiResource.AppConfigRequest
	if err := g.ShouldBindJSON(&data); err != nil {
		api.ResponseError(g, err)
		return
	}

	config := &appproject.AppConfig{
		Spec: appproject.AppConfigSpec{
			App: data.App,
			Config: data.Config,
		},
	}
	result, update, err := s.AppConfigService.Update(config)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"results": result, "update": update})
}
