package appservice

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appproject"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
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
		if _, ok := parentsMap[child.Spec.RootApp]; !ok{
			parentsMap[child.Spec.RootApp] = make([]*apiResource.Response, 0)
		}
		parentsMap[child.Spec.RootApp] = append(parentsMap[child.Spec.RootApp], child)
	}

	parentData, err := n.IService.ListByFilter(common.DefaultNamespace, common.AppProject, map[string]interface{}{"spec.parent_app":""}, sort, 0, 0)
	if err != nil {
		return nil, err
	}

	parents := make([]*apiResource.Response, 0)
	if err := utils.UnstructuredObjectToInstanceObj(parentData, &parents); err != nil {
		return nil, err
	}

	for i:=len(parents)-1;i>=0;i--{
		if _, ok := parentsMap[parents[i].UUID]; ok{
			parents[i].Children = parentsMap[parents[i].UUID]
		}else{
			parents = append(parents[:i], parents[i+1:]...)
		}
	}

	return parents, nil
}

func (n *NamespaceService) Create(request *apiResource.Request) (core.IObject, error) {

	req := &appproject.AppProject{
		Metadata: core.Metadata{
			Name: request.Name,
		},
		Spec: appproject.AppSpec{
			AppType:   appproject.Namespace,
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

	if err := n.IService.GetByFilter(common.DefaultNamespace, common.Namespace, req, filter); err == nil {
		return nil, errors.New("The Desc is exist")
	}

	req.GenerateVersion()
	parent := &appproject.AppProject{}
	if req.Spec.ParentApp != "" {
		if err := n.IService.GetByUUID(common.DefaultNamespace, common.AppProject, req.Spec.ParentApp, parent); err != nil {
			return nil, err
		}

		if parent.Spec.RootApp != "" {
			req.Spec.RootApp = parent.Spec.RootApp
		} else {
			req.Spec.RootApp = parent.Metadata.UUID
		}
	}

	return n.IService.Create(common.DefaultNamespace, common.Namespace, req)
}