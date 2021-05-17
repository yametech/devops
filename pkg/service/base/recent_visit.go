package base

import (
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
)

type RecentVisit struct {
	service.IService
}

func NewRecentVisit(i service.IService) *RecentVisit {
	return &RecentVisit{i}
}

func (r *RecentVisit) ListRecent(user string, page, pageSize int64) ([]*base.Module, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	//if user != "" {
	//	filter["spec.User"] = user
	//}
	filter["spec.User"] = user
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
		reverse(privateModule.Spec.Modules)
		for _, v := range privateModule.Spec.Modules {
			module := &base.Module{}
			err := r.IService.GetByUUID(common.DefaultNamespace, common.RecentVisit, v, module)
			if err != nil {
				return nil, err
			}
			moduleSlice = append(moduleSlice, module)
			return moduleSlice, nil
		}
	}
	return nil, errors.New("该用户没有最近访问记录！")
}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
