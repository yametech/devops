package base

import (
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/service"
)

type ShowAllGroupModule struct {
	service.IService
}

func NewShowAllGroupModule(i service.IService) *ShowAllGroupModule {
	return &ShowAllGroupModule{i}
}

func (s *ShowAllGroupModule) ListAllGroup(page, pageSize int64) ([]interface{}, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	filter["spec.parent"] = ""
	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}
	data, err := s.IService.ListByFilter(common.DefaultNamespace, common.AllModule, filter, sort, offset, pageSize)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *ShowAllGroupModule) ListAllModule(uuid string, page, pageSize int64) ([]interface{}, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	filter["spec.parent"] = uuid
	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}
	data, err := s.IService.ListByFilter(common.DefaultNamespace, common.AllModule, filter, sort, offset, pageSize)
	if err != nil {
		return nil, err
	}
	return data, nil
}
