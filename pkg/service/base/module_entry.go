package base

import (
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
)

type ModuleEntry struct {
	service.IService
}

func NewModuleEntry(i service.IService) *ModuleEntry {
	return &ModuleEntry{i}
}

func (m *ModuleEntry) CreateEntry(user, uuid string, page, pageSize int64) ([]*base.Module, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	//if user != "" {
	//	filter["spec.User"] = user
	//}
	filter["spec.user"] = user
	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}

	data, err := m.IService.ListByFilter(common.DefaultNamespace, common.ModuleEntry, filter, sort, offset, pageSize)
	if err != nil {
		return nil, err
	}
	modulates := &base.Module{}
	if err := m.IService.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, modulates); err != nil {
		return nil, errors.New("此模块的uuid在数据库中不存在！")
	}
	if data != nil {
		privateModule := &base.PrivateModule{}
		for _, v := range data {
			err := utils.UnstructuredObjectToInstanceObj(v, privateModule)
			if err != nil {
				return nil, err
			}
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
			err := m.IService.GetByUUID(common.DefaultNamespace, common.AllModule, v, module)
			if err != nil {
				return nil, err
			}
			moduleSlice = append(moduleSlice, module)
		}
		return moduleSlice, nil
	}

	uuidSlice := make([]string, 0)
	uuidSlice = append(uuidSlice, uuid)
	me := &base.PrivateModule{
		Metadata: core.Metadata{},
		Spec: base.PrivateModuleSpec{
			User:    user,
			Modules: uuidSlice,
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

func (m *ModuleEntry) DeleteEntry(user, uuid string, page, pageSize int64) ([]*base.Module, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	//if user != "" {
	//	filter["spec.User"] =user
	//}
	filter["spec.user"] = user
	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}

	data, err := m.IService.ListByFilter(common.DefaultNamespace, common.ModuleEntry, filter, sort, offset, pageSize)
	if err != nil {
		return nil, err
	}
	modulates := &base.Module{}
	if err := m.IService.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, modulates); err != nil {
		return nil, errors.New("此模块的uuid在数据库中不存在")
	}
	if data != nil {
		privateModule := &base.PrivateModule{}
		for _, v := range data {
			err := utils.UnstructuredObjectToInstanceObj(v, privateModule)
			if err != nil {
				return nil, err
			}
		}
		for k, v := range privateModule.Spec.Modules {
			if v == uuid {
				kk := k + 1
				privateModule.Spec.Modules = append(privateModule.Spec.Modules[:k], privateModule.Spec.Modules[kk:]...)
			}
		}
		_, judge, err := m.IService.Apply(common.DefaultNamespace, common.ModuleEntry, privateModule.UUID, privateModule, true)
		if !judge && err != nil {
			return nil, err
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
	return nil, errors.New("删除错误！")
}

func (m *ModuleEntry) QueryEntry(user string, page, pageSize int64) ([]*base.Module, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	//if user != "" {
	//	filter["spec.User"] =user
	//}
	filter["spec.user"] = user
	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}

	data, err := m.IService.ListByFilter(common.DefaultNamespace, common.ModuleEntry, filter, sort, offset, pageSize)
	if err != nil {
		return nil, err
	}
	if data != nil {
		privateModule := &base.PrivateModule{}
		for _, v := range data {
			err := utils.UnstructuredObjectToInstanceObj(v, privateModule)
			if err != nil {
				return nil, err
			}
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
	return nil, errors.New("查询错误,该用户没有添加快捷入口！")
}
