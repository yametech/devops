package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/appservice"
	"github.com/yametech/devops/pkg/resource/appservice"
	"strconv"
)

func (s *Server) GetAppConfig(g *gin.Context) {
	uuid := g.Param("uuid")
	config, err := s.AppConfigService.GetAppConfig(uuid)
	if err != nil {
		api.ResponseSuccess(g, &appservice.AppConfig{}, "")
		return
	}

	api.ResponseSuccess(g, config, "")
}

func (s *Server) GetAppResource(g *gin.Context) {
	uuid := g.Param("uuid")
	AppResource, err := s.AppConfigService.GetAppResources(uuid)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, AppResource, "")
}

func (s *Server) UpdateAppConfig(g *gin.Context) {
	data := &apiResource.AppConfigRequest{}
	if err := g.ShouldBindJSON(&data); err != nil {
		api.ResponseError(g, err)
		return
	}

	result, update, err := s.AppConfigService.UpdateAppConfig(data)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"results": result, "update": update}, "")
}

func (s *Server) UpdateAppResource(g *gin.Context) {
	data := &apiResource.NamespaceRequest{}
	if err := g.ShouldBindJSON(&data); err != nil {
		api.ResponseError(g, err)
		return
	}

	result, update, err := s.AppConfigService.UpdateConfigResource(data)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"result": result, "update": update}, "")
}

func (s *Server) DeleteResource(g *gin.Context) {
	uuid := g.Param("uuid")
	if err := s.AppConfigService.DeleteResource(uuid); err != nil {
		api.ResponseSuccess(g, gin.H{"delete": false}, "")
	}

	api.ResponseSuccess(g, gin.H{"delete": true}, "")
}

func (s *Server) ConfigHistory(g *gin.Context) {
	uuid := g.Param("uuid")
	page, err := strconv.ParseInt(g.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		api.ResponseError(g, errors.New("page need int type"))
		return
	}
	pageSize, err := strconv.ParseInt(g.DefaultQuery("page_size", "10"), 10, 64)
	if err != nil {
		api.ResponseError(g, errors.New("pageSize need int type"))
		return
	}

	results, err := s.AppConfigService.History(uuid, page, pageSize)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, results, "")
}
