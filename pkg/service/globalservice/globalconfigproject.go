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

func (a *GlobalConfigService) GetByUUID(name string, uuid string) (interface{}, error) {
	allies := &globalconfig.GlobalConfig{}
	if uuid != "" {
		if err := a.IService.GetByUUID(common.DefaultNamespace, common.GlobalConfig, uuid, allies); err != nil {
			return nil, err
		}
	}
	if name != "" {
		if val, ok := allies.Spec.Service[name]; ok {
			return val, nil
		}
	}
	return allies, nil
}

func (a *GlobalConfigService) Create(reqAll *globalconfigproject.RequestGlobalConfig) error {
	autoconfigure := &globalconfig.GlobalConfig{
		Metadata: core.Metadata{
			Name: reqAll.Name,
			Kind: reqAll.Kind,
		},
		Spec: globalconfig.Spec{
			Service: reqAll.Request.Service,
		},
	}
	autoconfigure.GenerateVersion()
	_, err := a.IService.Create(common.DefaultNamespace, common.GlobalConfig, autoconfigure)
	return err
}

func (a *GlobalConfigService) Update(uuid string, reqAll *globalconfigproject.RequestGlobalConfig) (core.IObject, bool, error) {
	autoconfigure := &globalconfig.GlobalConfig{
		Metadata: core.Metadata{
			Name: reqAll.Name,
			Kind: reqAll.Kind,
		},
		Spec: globalconfig.Spec{
			Service: reqAll.Request.Service,
		},
	}
	autoconfigure.GenerateVersion()
	return a.IService.Apply(common.DefaultNamespace, common.GlobalConfig, uuid, autoconfigure)
}
