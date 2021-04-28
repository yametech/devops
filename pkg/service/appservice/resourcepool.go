package appservice

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appservice"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResourcePoolService struct {
	service.IService
}

func NewResourcePoolService(i service.IService) *ResourcePoolService {
	return &ResourcePoolService{i}
}

func (n *ResourcePoolService) List() ([]*apiResource.Response, error) {
	parentsMap := make(map[string][]*apiResource.Response, 0)
	sort := map[string]interface{}{
		"metadata.create_time": 1,
	}

	childrenData, err := n.IService.ListByFilter(common.DefaultNamespace, common.ResourcePool, nil, sort, 0, 0)
	if err != nil {
		return nil, err
	}

	children := make([]*apiResource.Response, 0)
	if err = utils.UnstructuredObjectToInstanceObj(childrenData, &children); err != nil {
		return nil, err
	}

	for _, child := range children {
		if _, ok := parentsMap[child.Spec.RootApp]; !ok {
			parentsMap[child.Spec.RootApp] = make([]*apiResource.Response, 0)
		}
		parentsMap[child.Spec.RootApp] = append(parentsMap[child.Spec.RootApp], child)
	}

	parentData, err := n.IService.ListByFilter(common.DefaultNamespace, common.AppProject, map[string]interface{}{"spec.parent_app": ""}, sort, 0, 0)
	if err != nil {
		return nil, err
	}

	parents := make([]*apiResource.Response, 0)
	if err := utils.UnstructuredObjectToInstanceObj(parentData, &parents); err != nil {
		return nil, err
	}

	for i := len(parents) - 1; i >= 0; i-- {
		if _, ok := parentsMap[parents[i].UUID]; ok {
			parents[i].Children = parentsMap[parents[i].UUID]
		} else {
			parents = append(parents[:i], parents[i+1:]...)
		}
	}

	return parents, nil
}

func (n *ResourcePoolService) ListByLevel(level int, search string) (interface{}, error) {

	levelData := []func(int, string) (interface{}, error){
		n.ListAppProjectLevel,
		n.ListAppProjectLevel,
		n.ListAppProjectLevel,
		n.ListResourcePoolLevel,
	}

	if level > len(levelData)-1 {
		return nil, errors.New("have no this level")
	}

	return levelData[level](level, search)
}

func (n *ResourcePoolService) ListAppProjectLevel(level int, search string) (interface{}, error) {
	filter := make(map[string]interface{})
	filter["$or"] = []map[string]interface{}{
		{
			"spec.app_type": level,
			"metadata.name": bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
		},
		{
			"spec.app_type": level,
			"spec.desc":     bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
		},
	}

	sort := map[string]interface{}{
		"metadata.created_time": 1,
	}

	return n.IService.ListByFilter(common.DefaultNamespace, common.AppProject, filter, sort, 0, 0)
}

func (n *ResourcePoolService) ListResourcePoolLevel(level int, search string) (interface{}, error) {
	filter := make(map[string]interface{})
	filter["$or"] = []map[string]interface{}{
		{
			"metadata.name": bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
		},
		{
			"spec.desc": bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
		},
	}

	sort := map[string]interface{}{
		"metadata.created_time": 1,
	}
	return n.IService.ListByFilter(common.DefaultNamespace, common.ResourcePool, filter, sort, 0, 0)
}

func (n *ResourcePoolService) Create(request *apiResource.Request) (core.IObject, error) {

	req := &appservice.ResourcePool{
		Metadata: core.Metadata{
			Name: request.Name,
		},
		Spec: appservice.ResourcePoolSpec{
			ParentApp: request.ParentApp,
			Desc:      request.Desc,
		},
	}

	if req.Spec.Desc == "" {
		return nil, errors.New("The Desc is requried")
	}

	filter := map[string]interface{}{
		"spec.desc": req.Spec.Desc,
	}

	if err := n.IService.GetByFilter(common.DefaultNamespace, common.ResourcePool, req, filter); err == nil {
		return nil, errors.New("The Desc is exist")
	}

	req.GenerateVersion()
	parent := &appservice.AppProject{}
	if req.Spec.ParentApp != "" {
		if err := n.IService.GetByUUID(common.DefaultNamespace, common.AppProject, req.Spec.ParentApp, parent); err != nil {
			return nil, err
		}
		req.Spec.RootApp = parent.Metadata.UUID
	}

	return n.IService.Create(common.DefaultNamespace, common.ResourcePool, req)
}
