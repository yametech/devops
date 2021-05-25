package base

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/base"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChildModuleService struct {
	service.IService
}

func NewChildModuleService(i service.IService) *ChildModuleService {
	return &ChildModuleService{IService: i}
}

func (g *ChildModuleService) CreateChildModule(request *apiResource.ModuleRequest) (core.IObject, error){
	if err := g.IService.GetByFilter(common.DefaultNamespace, common.ChildModule, base.Module{}, map[string]interface{}{
		"metadata.name": request.Name,
		"spec.parent": request.Parent,
	}); err == nil {
		return nil, errors.New("The ChildModule Name is exists")
	}

	req := &base.Module{
		Metadata: core.Metadata{
			Name: request.Name,
		},
		Spec: base.ModuleSpec{
			Parent: request.Parent,
			Extends: request.Extends,
		},
	}

	req.GenerateVersion()
	return g.IService.Create(common.DefaultNamespace, common.ChildModule, req)
}

func (g *ChildModuleService) ListChildModule(parent, search string, page, pageSize int64) ([]interface{}, int64, error) {
	title := &base.Module{}
	if err := g.IService.GetByUUID(common.DefaultNamespace, common.AllModule, parent, title); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	filter["spec.parent"] = parent
	if search != "" {
		filter["metadata.name"] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}}
	}
	sort := map[string]interface{}{
		"metadata.created_time": 1,
	}

	data, err := g.IService.ListByFilter(common.DefaultNamespace, common.ChildModule, filter, sort, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := g.IService.Count(common.DefaultNamespace, common.ChildModule, filter)
	if err != nil {
		return nil, 0, err
	}

	result := make([]interface{}, 0)
	result = append(result, title)
	result = append(result, data...)

	return result, count, nil
}

func (g *ChildModuleService) DeleteChildModule(uuid string) (bool, error) {
	if err := g.IService.Delete(common.DefaultNamespace, common.ChildModule, uuid); err != nil {
		return false, err
	}
	return true, nil
}