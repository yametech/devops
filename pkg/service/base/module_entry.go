package base

import (
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
)

type ModuleEntry struct {
	service.IService
}

func NewModuleEntry(i service.IService) *ModuleEntry {
	return &ModuleEntry{i}
}

func (m *ModuleEntry) CreateEntry(user, uuid string) ([]*base.Module, error) {
	filter := map[string]interface{}{
		"spec.user": user,
	}
	privateModule := &base.PrivateModule{}
	if err := m.IService.GetByFilter(common.DefaultNamespace, common.ModuleEntry, privateModule, filter); err != nil {
		me := &base.PrivateModule{
			Metadata: core.Metadata{},
			Spec: base.PrivateModuleSpec{
				User:    user,
				Modules: []string{uuid},
			},
		}
		_, err = m.IService.Create(common.DefaultNamespace, common.ModuleEntry, me)
		if err != nil {
			return nil, err
		}
		module := &base.Module{}
		moduleSlice := make([]*base.Module, 0)
		err = m.IService.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, module)
		if err != nil {
			return nil, err
		}
		moduleSlice = append(moduleSlice, module)
		return moduleSlice, nil
	}
	modulates := &base.Module{}
	if err := m.IService.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, modulates); err != nil {
		return nil, errors.New("此模块的uuid在数据库中不存在！")
	}
	for _, v := range privateModule.Spec.Modules {
		if v == uuid {
			return nil, errors.New("此模块快捷入口已经存在！")
		}
	}
	privateModule.Spec.Modules = append(privateModule.Spec.Modules, uuid)
	_, judge, err := m.IService.Apply(common.DefaultNamespace, common.ModuleEntry, privateModule.UUID, privateModule, true)
	if !judge && err != nil {
		return nil, err
	}
	moduleSlice := make([]*base.Module, 0)
	for _, v := range privateModule.Spec.Modules {
		module := &base.Module{}
		_err := m.IService.GetByUUID(common.DefaultNamespace, common.AllModule, v, module)
		if _err != nil {
			return nil, _err
		}
		moduleSlice = append(moduleSlice, module)
	}
	return moduleSlice, nil
}

func (m *ModuleEntry) DeleteEntry(user, uuid string) ([]*base.Module, error) {
	filter := map[string]interface{}{
		"spec.user": user,
	}

	privateModule := &base.PrivateModule{}
	if err := m.IService.GetByFilter(common.DefaultNamespace, common.ModuleEntry, privateModule, filter); err != nil {
		return nil, errors.New("此用户数据在数据库中不存在！")
	}
	modulates := &base.Module{}
	if err := m.IService.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, modulates); err != nil {
		return nil, errors.New("此模块的uuid在数据库中不存在！")
	}
	for k, v := range privateModule.Spec.Modules {
		if v == uuid {
			privateModule.Spec.Modules = append(privateModule.Spec.Modules[:k], privateModule.Spec.Modules[k+1:]...)
		}
	}
	_, judge, err := m.IService.Apply(common.DefaultNamespace, common.ModuleEntry, privateModule.UUID, privateModule, true)
	if !judge && err != nil {
		return nil, err
	}
	moduleSlice := make([]*base.Module, 0)
	for _, v := range privateModule.Spec.Modules {
		module := &base.Module{}
		_err := m.IService.GetByUUID(common.DefaultNamespace, common.AllModule, v, module)
		if _err != nil {
			return nil, _err
		}
		moduleSlice = append(moduleSlice, module)
	}
	return moduleSlice, nil
}

func (m *ModuleEntry) QueryEntry(user string) ([]*base.Module, error) {
	filter := map[string]interface{}{
		"spec.user": user,
	}
	privateModule := &base.PrivateModule{}
	if err := m.IService.GetByFilter(common.DefaultNamespace, common.ModuleEntry, privateModule, filter); err != nil {
		return nil, errors.New("此用户数据在数据库中不存在！")
	}
	moduleSlice := make([]*base.Module, 0)
	for _, v := range privateModule.Spec.Modules {
		module := &base.Module{}
		err := m.IService.GetByUUID(common.DefaultNamespace, common.AllModule, v, module)
		if err != nil {
			return nil, err
		}
		moduleSlice = append(moduleSlice, module)
	}
	return moduleSlice, nil
}
