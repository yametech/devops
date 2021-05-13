package module_entry

import (
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModuleEntry struct {
	service.IService
}

func NewModuleEntry(i service.IService) *ModuleEntry {
	return &ModuleEntry{i}
}

func (m *ModuleEntry) Create(user, uuid string, page, pageSize int64) ([]*base.Module, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	if user != "" {
		filter["spec.User"] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + user + ".*", Options: "i"}}
	}
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
			return moduleSlice, nil
		}
	} else {
		uuidSlice := make([]string, 0)
		uuidSlice = append(uuidSlice, uuid)
		me := &base.PrivateModule{
			Metadata: core.Metadata{},
			Spec: base.PrivateModuleSpec{
				User:    user,
				Modules: uuidSlice,
			},
		}

		_, err := m.IService.Create(common.DefaultNamespace, common.ModuleEntry, me)
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
	return nil, err
}

func (m *ModuleEntry) Delete(user, uuid string, page, pageSize int64) ([]*base.Module, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	if user != "" {
		filter["spec.User"] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + user + ".*", Options: "i"}}
	}
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
			return moduleSlice, nil
		}
	}
	return nil, err
}
