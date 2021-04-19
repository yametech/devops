package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/resource/appproject"
	"net/http"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
)

func (s *Server) GetAppConfig(g *gin.Context) {
	var data apiResource.AppConfigRequest
	if err := g.ShouldBindJSON(&data); err != nil {
		api.RequestParamsError(g, "error", err)
	}

	config := &appproject.AppConfig{
		Spec: appproject.AppConfigSpec{
			App: data.App,
			ConfigType: data.ConfigType,
		},
	}

	if err := s.AppConfigService.GetByFilter(config); err != nil {
		api.RequestParamsError(g, "error", err)
	}

	g.JSON(http.StatusOK, gin.H{
		"data": config,
	})
}

func (s *Server) UpdateAppConfig(g *gin.Context) {
	var data apiResource.AppConfigRequest
	if err := g.ShouldBindJSON(&data); err != nil {
		api.RequestParamsError(g, "error", err)
	}

	config := &appproject.AppConfig{
		Spec: appproject.AppConfigSpec{
			App: data.App,
			ConfigType: data.ConfigType,
			Config: data.Config,
		},
	}
	result, update, err := s.AppConfigService.Update(config)
	if err != nil {
		api.RequestParamsError(g, "error", err)
	}

	g.JSON(http.StatusOK, gin.H{
		"data":   result,
		"update": update,
	})
}
