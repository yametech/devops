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

type GlobalModuleService struct {
	service.IService
}

func NewGlobalModuleService(i service.IService) *GlobalModuleService {
	return &GlobalModuleService{IService: i}
}

func (g *GlobalModuleService) CreateGlobalModule(request *apiResource.ModuleRequest) (core.IObject, error){
	if err := g.IService.GetByFilter(common.DefaultNamespace, common.GlobalModule, base.Module{}, map[string]interface{}{
		"metadata.name": request.Name,
	}); err == nil {
		return nil, errors.New("The GlobalModule Name is exists")
	}

	req := &base.Module{
		Metadata: core.Metadata{
			Name: request.Name,
		},
		Spec: base.ModuleSpec{
			Extends: request.Extends,
		},
	}

	req.GenerateVersion()
	return g.IService.Create(common.DefaultNamespace, common.GlobalModule, req)
}

func (g *GlobalModuleService) ListGlobalModule(search string, page, pageSize int64) ([]interface{}, int64, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	if search != "" {
		filter["metadata.name"] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}}
	}
	sort := map[string]interface{}{
		"metadata.created_time": 1,
	}

	data, err := g.IService.ListByFilter(common.DefaultNamespace, common.GlobalModule, filter, sort, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := g.IService.Count(common.DefaultNamespace, common.GlobalModule, filter)
	if err != nil {
		return nil, 0, err
	}
	return data, count, nil
}

func (g *GlobalModuleService) DeleteGlobalModule(uuid string) (bool, error) {
	if err := g.IService.Delete(common.DefaultNamespace, common.GlobalModule, uuid); err != nil {
		return false, err
	}
	return true, nil
}