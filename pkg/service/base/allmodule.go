package base

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/base"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AllModuleService struct {
	service.IService
}

func NewAllModuleService(i service.IService) *AllModuleService {
	return &AllModuleService{IService: i}
}

func (a *AllModuleService) CreateGroup(req *apiResource.ModuleRequest) (core.IObject, error) {

	dbModule := &base.Module{}
	if err := a.IService.GetByFilter(common.DefaultNamespace, common.AllModule, dbModule, map[string]interface{}{
		"metadata.name": req.Name,
	}); err == nil {
		return nil, errors.New("The module name is exists")
	}

	module := &base.Module{
		Metadata: core.Metadata{
			Name: req.Name,
		},
		Spec: base.ModuleSpec{
			Extends: req.Extends,
		},
	}

	return a.IService.Create(common.DefaultNamespace, common.AllModule, module)
}

func (a *AllModuleService) CreateModule(req *apiResource.ModuleRequest) (core.IObject, error) {
	dbModule := &base.Module{}
	if err := a.IService.GetByFilter(common.DefaultNamespace, common.AllModule, dbModule, map[string]interface{}{
		"spec.parent":   req.Parent,
		"metadata.name": req.Name,
	}); err == nil {
		return nil, errors.New("The module name is exists")
	}

	module := &base.Module{
		Metadata: core.Metadata{
			Name: req.Name,
		},
		Spec: base.ModuleSpec{
			Parent:  req.Parent,
			Extends: req.Extends,
		},
	}

	return a.IService.Create(common.DefaultNamespace, common.AllModule, module)
}

func (a *AllModuleService) DeleteGroupAndModule(uuid string) (bool, error) {
	if err := a.IService.Delete(common.DefaultNamespace, common.AllModule, uuid); err != nil {
		return false, err
	}
	return true, nil
}

func (a *AllModuleService) ListAll(search string) ([]*apiResource.ModuleResponse, error) {
	filter := map[string]interface{}{
		"spec.parent": "",
	}

	sort := map[string]interface{}{
		"metadata.created_time": 1,
	}

	data, err := a.IService.ListByFilter(common.DefaultNamespace, common.AllModule, filter, sort, 0, 0)
	if err != nil {
		return nil, err
	}

	groups := make([]*apiResource.ModuleResponse, 0)
	if err = utils.UnstructuredObjectToInstanceObj(data, &groups); err != nil {
		return nil, err
	}

	for i := len(groups) - 1; i >= 0; i-- {
		data, err = a.IService.ListByFilter(common.DefaultNamespace, common.AllModule, map[string]interface{}{
			"spec.parent":   groups[i].UUID,
			"metadata.name": bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
		}, sort, 0, 0)
		if err != nil {
			return nil, err
		}
		modules := make([]*base.Module, 0)
		if err = utils.UnstructuredObjectToInstanceObj(data, &modules); err != nil {
			return nil, err
		}

		if len(modules) > 0 {
			groups[i].Children = modules
		} else {
			if search != "" {
				groups = append(groups[:i], groups[i+1:]...)
			}
		}
	}

	return groups, nil
}
