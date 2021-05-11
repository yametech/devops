package appservice

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appservice"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appservice"
	"github.com/yametech/devops/pkg/resource/workorder"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type NamespaceService struct {
	service.IService
}

func NewNamespaceService(i service.IService) *NamespaceService {
	return &NamespaceService{i}
}

func (n *NamespaceService) List() ([]*apiResource.Response, error) {
	parentsMap := make(map[string][]*apiResource.Response, 0)
	sort := map[string]interface{}{
		"metadata.create_time": 1,
	}

	childrenData, err := n.IService.ListByFilter(common.DefaultNamespace, common.Namespace, nil, sort, 0, 0)
	if err != nil {
		return nil, err
	}

	children := make([]*apiResource.Response, 0)
	if err = utils.UnstructuredObjectToInstanceObj(childrenData, &children); err != nil {
		return nil, err
	}

	for _, child := range children {
		if _, ok := parentsMap[child.Spec.ParentApp]; !ok {
			parentsMap[child.Spec.ParentApp] = make([]*apiResource.Response, 0)
		}
		parentsMap[child.Spec.ParentApp] = append(parentsMap[child.Spec.ParentApp], child)
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

func (n *NamespaceService) ListByLevel(level int, search string, filter string) (interface{}, error) {

	levelData := []func(int, string, string) (interface{}, error){
		n.ListAppProjectLevel,
		n.ListAppProjectLevel,
		n.ListAppProjectLevel,
		n.ListResourcePoolLevel,
	}

	if level > len(levelData)-1 {
		return nil, errors.New("have no this level")
	}

	return levelData[level](level, search, filter)
}

func (n *NamespaceService) ListAppProjectLevel(level int, search string, parentApp string) (interface{}, error) {
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

func (n *NamespaceService) ListResourcePoolLevel(level int, search string, parentApp string) (interface{}, error) {
	filter := make(map[string]interface{})
	if parentApp != ""{
		filter["$or"] = []map[string]interface{}{
			{
				"metadata.name": bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
				"spec.parent_app": parentApp,
			},
			{
				"spec.desc": bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
				"spec.parent_app": parentApp,
			},
		}
	}else{
		filter["$or"] = []map[string]interface{}{
			{
				"metadata.name": bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
			},
			{
				"spec.desc": bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
			},
		}
	}

	sort := map[string]interface{}{
		"metadata.created_time": 1,
	}
	return n.IService.ListByFilter(common.DefaultNamespace, common.Namespace, filter, sort, 0, 0)
}

func (n *NamespaceService) Update(request *apiResource.Request) (core.IObject, error) {

	obj := &appservice.Namespace{}
	if err := n.IService.GetByUUID(common.DefaultNamespace, common.Namespace, request.UUID, obj); err != nil {
		log.Printf("namespace update error: %v\n", err)
	}

	if request.Desc == "" {
		return nil, errors.New("The Desc is requried")
	}

	filter := map[string]interface{}{
		"spec.desc": request.Desc,
		"metadata.uuid": bson.M{"$ne": request.UUID},
	}

	req := &appservice.Namespace{}
	if err := n.IService.GetByFilter(common.DefaultNamespace, common.Namespace, req, filter); err == nil {
		return nil, errors.New("The Desc is exist")
	}

	history := &appservice.NamespaceHistory{}
	history.Spec.Creator = ""

	// get from cmdb
	cpus := 1000.0
	memory := 10024000

	history.Spec.Before = map[string]interface{}{
		"cpu": cpus,
		"memory": memory,
	}

	obj.Name = request.Name
	obj.Spec.ParentApp = request.ParentApp
	obj.Spec.Desc = request.Desc

	req.GenerateVersion()
	result, _, err := n.IService.Apply(common.DefaultNamespace, common.Namespace, obj.UUID, obj, true)
	if err != nil {
		return nil, err
	}

	history.Spec.App = obj.UUID
	history.Spec.Now = map[string]interface{}{
		"cpu": request.Cpu,
		"memory": request.Memory,
	}

	if _, err = n.IService.Create(common.DefaultNamespace, common.History, history); err != nil {
		return nil, err
	}

	return result, nil
}

func (n *NamespaceService) OrderToNamespaceSuccess(obj *workorder.WorkOrder) error {

	filter := map[string]interface{}{
		"spec.parent_app": obj.UUID,
	}

	count, err := n.IService.Count(common.DefaultNamespace, common.Namespace, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("the OrderToNamespaceCheck namespace is exist")
	}

	configs := make([]*apiResource.Request, 0)
	if err = utils.UnstructuredObjectToInstanceObj(obj.Spec.Extends, &configs); err != nil {
		return err
	}

	for _, config := range configs {
		if _, err = n.Update(config); err != nil {
			return err
		}
	}

	return nil
}
