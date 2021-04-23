package globalservice

import (
	"github.com/yametech/devops/pkg/api/resource/globalconfigproject"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/globalconfig"
	"github.com/yametech/devops/pkg/service"
)

type GlobalConfigService struct {
	service.IService
}

func NewAllConfigService(i service.IService) *GlobalConfigService {
	return &GlobalConfigService{i}
}

func (a *GlobalConfigService) List(page, pageSize int64) ([]interface{}, error) {
	offset := (page - 1) * pageSize
	sort := map[string]interface{}{
		"metadata.version": -1,
	}

	unStruct, err := a.IService.List(common.DefaultNamespace, common.GlobalConfig, "", sort, offset, pageSize)

	return unStruct, err
}

func (a *GlobalConfigService) Create(reqAll *globalconfigproject.RequestGlobalConfig) error {
	autoconfigure := &globalconfig.GlobalConfig{
		Spec: globalconfig.Spec{
			Service: reqAll.Service,
		},
	}
	autoconfigure.GenerateVersion()
	_, err := a.IService.Create(common.DefaultNamespace, common.GlobalConfig, autoconfigure)
	return err
}

func (a *GlobalConfigService) Update(uuid string, reqAll *globalconfigproject.RequestGlobalConfig, forceApply bool) (core.IObject, bool, error) {
	autoconfigure := &globalconfig.GlobalConfig{
		Spec: globalconfig.Spec{
			Service: reqAll.Service,
		},
	}

	autoconfigure.GenerateVersion()
	updateObject, whether, err := a.IService.Apply(common.DefaultNamespace, common.GlobalConfig, uuid, autoconfigure, forceApply)
	if autoconfigure.Name == "" {
		autoconfigure.Name = "全局配置服务"
	}
	return updateObject, whether, err
}
