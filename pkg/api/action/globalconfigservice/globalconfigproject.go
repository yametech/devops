package globalconfigservice

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/api/resource/globalconfigproject"
)

func (s *Server) ListGlobalConfig(c *gin.Context) {
	pageInt := cast.ToUint(c.DefaultQuery("page", "1"))
	pageSizeInt := cast.ToUint(c.DefaultQuery("pageSize", "10"))
	res, err := s.GlobalConfigService.List(int64(pageInt), int64(pageSizeInt))
	if err != nil {
		api.ResponseError(c, err)
		return
	}
	api.ResponseSuccess(c, res, "")
}

func (s *Server) CreateGlobalConfig(c *gin.Context) {
	request := &globalconfigproject.RequestGlobalConfig{}
	if err := c.ShouldBindJSON(&request); err != nil {
		api.ResponseError(c, err)
		return
	}
	res, err := s.GlobalConfigService.Create(request)
	if err != nil {
		api.ResponseError(c, err)
		return
	}
	api.ResponseSuccess(c, res, "")
}

func (s *Server) UpdateGlobalConfig(c *gin.Context) {
	request := &globalconfigproject.RequestGlobalConfig{}
	if err := c.ShouldBindJSON(&request); err != nil {
		api.ResponseError(c, err)
		return
	}
	data, _, err := s.GlobalConfigService.Update(globalconfigproject.RequestGlobalConfigUUID, request)
	if err != nil {
		api.ResponseError(c, err)
		return
	}
	api.ResponseSuccess(c, data, "")
}
