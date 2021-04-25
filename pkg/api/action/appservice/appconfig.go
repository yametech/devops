package appservice

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"github.com/yametech/devops/pkg/resource/appproject"
	"strconv"
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

	api.ResponseSuccess(g, results)
}
