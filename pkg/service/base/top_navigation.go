package base

import (
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
)

type TopNavigation struct {
	service.IService
}

func NewTopNavigation(i service.IService) *TopNavigation {
	return &TopNavigation{i}
}

func (t *TopNavigation) ListTopModule(user string) ([]interface{}, error) {
	TopModule := &base.PrivateModule{}
	response := make([]interface{}, 0)
	if err := t.GetByFilter(common.DefaultNamespace, common.AllModule, TopModule, map[string]interface{}{"spec.user": user}); err != nil {
		return nil, errors.New("在数据库中找不到此用户数据！")
	}
	for _, uuid := range TopModule.Spec.Modules {
		data := &base.Module{}
		if err := t.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, data); err != nil {
			return nil, err
		}
		response = append(response, data)
	}
	return response, nil
}

func (t *TopNavigation) CreateTopModule(uuid, user string) (core.IObject, bool, error) {
	TopModule := &base.PrivateModule{}
	module := &base.Module{}
	if err := t.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, module); err != nil {
		return nil, false, errors.New("在数据库中找不到此uuid的相应模块！")
	}
	if err := t.GetByFilter(common.DefaultNamespace, common.Topmodule, TopModule, map[string]interface{}{"spec.user": user}); err != nil {
		TopModule.Spec.Modules = make([]string, 0)
		TopModule.Spec.User = user
	}
	TopModule.Spec.Modules = append(TopModule.Spec.Modules, uuid)
	return t.Apply(common.DefaultNamespace, common.Topmodule, TopModule.UUID, TopModule, true)
}

func (t *TopNavigation) DeleteTopModule(uuid, user string) (core.IObject, bool, error) {
	TopModule := &base.PrivateModule{}
	module := &base.Module{}
	if err := t.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, module); err != nil {
		return nil, false, errors.New("在数据库中找不到此uuid的相应模块！")
	}
	if err := t.GetByFilter(common.DefaultNamespace, common.Topmodule, TopModule, map[string]interface{}{"spec.user": user}); err != nil {
		return nil, false, errors.New("在数据库中找不到此用户数据！")
	}
	for i := len(TopModule.Spec.Modules) - 1; i > 0; i-- {
		if TopModule.Spec.Modules[i] == uuid {
			TopModule.Spec.Modules = append(TopModule.Spec.Modules[:i], TopModule.Spec.Modules[i+1:]...)
		}
	}
	return t.Apply(common.DefaultNamespace, common.Topmodule, TopModule.UUID, TopModule, true)
}
