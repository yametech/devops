package base

import (
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/store"
	"github.com/yametech/devops/pkg/utils"
)

type RecentVisit struct {
	service.IService
}

var _ store.IKVStore = (*RecentVisit)(nil)
var _ service.IService = (*RecentVisit)(nil)

func NewRecentVisit(i service.IService) *RecentVisit {
	return &RecentVisit{i}
}

func (r *RecentVisit) CreateRecent(user, uuid string, page, pageSize int64) ([]*base.Module, error) {
	modulates := &base.Module{}
	if err := r.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, modulates); err != nil {
		return nil, errors.New("此模块的uuid在数据库中不存在！")
	}
	privateModule := &base.PrivateModule{}
	err := r.GetByFilter(common.DefaultNamespace, common.RecentVisit, privateModule, map[string]interface{}{"spec.user": user})
	if err != nil {
		me := &base.PrivateModule{
			Metadata: core.Metadata{},
			Spec: base.PrivateModuleSpec{
				User:    user,
				Modules: []string{uuid},
			},
		}
		_, err = r.IService.Create(common.DefaultNamespace, common.RecentVisit, me)
		if err != nil {
			return nil, err
		}
		module := &base.Module{}
		moduleSlice := make([]*base.Module, 0)
		err = r.IService.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, module)
		if err != nil {
			return nil, err
		}
		moduleSlice = append(moduleSlice, module)
		return moduleSlice, nil
	}
	if len(privateModule.Spec.Modules) < 6 {
		privateModule.Spec.Modules = append(privateModule.Spec.Modules, uuid)
		privateModule.Spec.User = user
		_, judge, _err := r.Apply(common.DefaultNamespace, common.RecentVisit, privateModule.UUID, privateModule, true)
		if !judge && _err != nil {
			return nil, errors.New("最近访问更新失败！")
		}
	} else {
		privateModule.Spec.Modules = append(privateModule.Spec.Modules[1:], uuid)
		_, judge, _err := r.Apply(common.DefaultNamespace, common.RecentVisit, privateModule.UUID, privateModule, true)
		if !judge && _err != nil {
			return nil, errors.New("最近访问更新失败！")
		}
	}
	return r.ListRecent(user, page, pageSize)
}

func (r *RecentVisit) ListRecent(user string, page, pageSize int64) ([]*base.Module, error) {
	offset := (page - 1) * pageSize
	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}
	data, err := r.IService.ListByFilter(common.DefaultNamespace, common.RecentVisit, map[string]interface{}{"spec.user": user}, sort, offset, pageSize)
	if err != nil {
		return nil, err
	}
	if data != nil {
		privateModule := &base.PrivateModule{}
		for _, v := range data {
			_err := utils.UnstructuredObjectToInstanceObj(v, privateModule)
			if _err != nil {
				return nil, _err
			}
		}
		moduleSlice := make([]*base.Module, 0)
		for i := len(privateModule.Spec.Modules) - 1; i >= 0; i-- {
			module := &base.Module{}
			_err := r.IService.GetByUUID(common.DefaultNamespace, common.AllModule, privateModule.Spec.Modules[i], module)
			if _err != nil {
				return nil, _err
			}
			moduleSlice = append(moduleSlice, module)
		}
		return moduleSlice, nil
	}
	return nil, errors.New("该用户没有最近访问记录！")
}
