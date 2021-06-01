package base

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/common"
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
	privateModule := &base.PrivateModule{}
	err := r.GetByFilter(common.DefaultNamespace, common.RecentVisit, privateModule, map[string]interface{}{"spec.user": user})
	if err != nil {
		return nil, errors.New("没有找到当前用户！")
	}
	modulates := &base.Module{}
	if err := r.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, modulates); err != nil {
		return nil, errors.New("此模块的uuid在数据库中不存在！")
	}
	if len(privateModule.Spec.Modules) < 6 {
		privateModule.Spec.Modules = append(privateModule.Spec.Modules, uuid)
		privateModule.Spec.User = user
		_, judge, err := r.Apply(common.DefaultNamespace, common.RecentVisit, privateModule.UUID, privateModule, true)
		if !judge && err != nil {
			return nil, errors.New("最近访问更新失败！")
		}
	} else {
		fmt.Printf("%p", privateModule)
		privateModule.Spec.Modules = append(privateModule.Spec.Modules[1:], uuid)
		fmt.Printf("%p", privateModule)
		//privateModule.Spec.Modules = append([]string{uuid},privateModule.Spec.Modules[1:]...)
		//
		//privateModule.Spec.Modules = append(privateModule.Spec.Modules[1:len(privateModule.Spec.Modules):len(privateModule.Spec.Modules)-1],uuid)

		_, judge, err := r.Apply(common.DefaultNamespace, common.RecentVisit, privateModule.UUID, privateModule, true)
		if !judge && err != nil {
			return nil, errors.New("最近访问更新失败！")
		}
	}
	return r.ListRecent(user, page, pageSize)
}

func (r *RecentVisit) ListRecent(user string, page, pageSize int64) ([]*base.Module, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	//if user != "" {
	//	filter["spec.User"] = user
	//}
	filter["spec.user"] = user
	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}

	data, err := r.IService.ListByFilter(common.DefaultNamespace, common.RecentVisit, filter, sort, offset, pageSize)
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
		for i := len(privateModule.Spec.Modules) - 1; i >= 0; i-- {
			module := &base.Module{}
			err := r.IService.GetByUUID(common.DefaultNamespace, common.AllModule, privateModule.Spec.Modules[i], module)
			if err != nil {
				return nil, err
			}
			moduleSlice = append(moduleSlice, module)
		}
		return moduleSlice, nil
	}
	return nil, errors.New("该用户没有最近访问记录！")
}
