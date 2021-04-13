package appservice

import (
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/resource"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
)

type AppProjectService struct {
	service.IService
}

func NewAppProjectService(i service.IService) *AppProjectService {
	return &AppProjectService{i}
}

func (a *AppProjectService) Children(req *resource.AppProjectResponse, sort map[string]interface{}) error {
	filter := make(map[string]interface{}, 0)
	filter["spec.parent_app"] = req.UUID
	data, count, err := a.IService.ListByFilter(common.DefaultNamespace, common.AppProject, filter, sort, 0, 0)
	if err != nil{
		return err
	}

	if count > 0{
		children := make([]*resource.AppProjectResponse, 0)
		err = utils.Clone(data, &children)
		if err != nil{
			return err
		}

		for _, child := range children{
			_child := child
			err := a.Children(_child, sort)
			if err != nil{
				return err
			}
		}

		req.Children = children
	}

	return nil
}

func (a *AppProjectService) List(page, pageSize int64) ([]interface{}, int64, error) {
	offset := (page - 1) * pageSize
	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}

	// Get the BusinessLine
	businessLine := make([]*resource.AppProjectResponse, 0)
	filter := map[string]interface{}{
		"spec.parent_app": "",
	}
	data, count, err := a.IService.ListByFilter(common.DefaultNamespace, common.AppProject, filter, sort, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// To []*resource.AppProjectResponse
	err = utils.Clone(data, &businessLine)
	if err != nil{
		return nil, 0, err
	}

	// Get the Children of BusinessLine
	for _, line := range businessLine{
		_line := line
		err := a.Children(_line, sort)
		if err != nil{
			return nil, 0, err
		}
	}

	// To []interface{}
	res := make([]interface{}, 0)
	err = utils.Clone(businessLine, &res)
	if err != nil{
		return nil, 0, err
	}

	return res, count, nil
}

func (a *AppProjectService) Create(req *resource.AppProject) error {
	req.GenerateVersion()
	_, err := a.IService.Create(common.DefaultNamespace, common.AppProject, req)
	if err != nil {
		return err
	}
	return nil
}
