package globalconfigservice

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/api/resource/globalconfigproject"
	"strconv"
)

func (s *Server) ListGlobalConfig(g *gin.Context) {
	pageInt, _ := strconv.Atoi(g.DefaultQuery("page", "1"))
	pageSizeInt, _ := strconv.Atoi(g.DefaultQuery("pageSize", "10"))
	res, err := s.GlobalConfigService.List(int64(pageInt), int64(pageSizeInt))
	if err != nil {
		api.ResponseError(g, err)
		return
	}
	api.ResponseSuccess(g, res, "")
}

func (s *Server) CreateGlobalConfig(g *gin.Context) {
	rawData, err := g.GetRawData()
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}
	request := &globalconfigproject.RequestGlobalConfig{}
	if err := json.Unmarshal(rawData, &request); err != nil {
		api.RequestParamsError(g, "unmarshal json error", err)
		return
	}
	res, err := s.GlobalConfigService.Create(request)
	if err != nil {
		api.ResponseError(g, err)
		return
	}
	api.ResponseSuccess(g, res, "")
}

func (s *Server) UpdateGlobalConfig(g *gin.Context) {
	rawData, err := g.GetRawData()
	if err != nil {
		api.RequestParamsError(g, "get rawData error", err)
		return
	}
	request := &globalconfigproject.RequestGlobalConfig{}
	if err := json.Unmarshal(rawData, &request); err != nil {
		api.RequestParamsError(g, "unmarshal json error", err)
		return
	}
	data, _, err := s.GlobalConfigService.Update(globalconfigproject.RequestGlobalConfigUUID, request)
	if err != nil {
		api.ResponseError(g, err)
		return
	}
	api.ResponseSuccess(g, data, "")
}
