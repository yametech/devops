package appservice

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/apppservice"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appservice"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	"log"
)

type ResourcePoolConfigService struct {
	service.IService
}

func NewResourcePoolConfigService(i service.IService) *ResourcePoolConfigService {
	return &ResourcePoolConfigService{IService: i}
}

func (n *ResourcePoolConfigService) GetResourcePoolConfig(appid string) (core.IObject, error) {
	req := &appservice.AppResource{
		Spec: appservice.AppResourceSpec{
			App: appid,
		},
	}

	if err := n.IService.GetByFilter(common.DefaultNamespace, common.AppResource, req, map[string]interface{}{
		"spec.app": req.Spec.App,
	}); err != nil {
		return nil, err
	}

	return req, nil
}

func (n *ResourcePoolConfigService) Update(data *apiResource.NamespaceRequest) (core.IObject, bool, error) {

	namespace := &appservice.Namespace{}
	if err := n.GetByUUID(common.DefaultNamespace, common.Namespace, data.App, namespace); err != nil {
		return nil, false, errors.New("The Namespace is not exist")
	}

	dbObj := &appservice.AppResource{}
	err := n.IService.GetByFilter(common.DefaultNamespace, common.AppResource, dbObj, map[string]interface{}{
		"spec.app": namespace.Metadata.UUID,
	})
	if err != nil {
		log.Printf("Update AppResource Not Found Create New One: %v\n", err)
	}


	// create history
	// Get creator
	history := &appservice.AppResourceHistory{}
	history.Spec.Creator = ""
	history.Spec.Before = dbObj
	history.Spec.App = namespace.Metadata.UUID

	dbObj.Spec.App = namespace.Metadata.UUID
	dbObj.Spec.Threshold = data.Threshold
	dbObj.Spec.Approval = data.Approval
	dbObj.Spec.Cpu = data.Cpu
	dbObj.Spec.Memory = data.Memory
	dbObj.Spec.Pod = data.Pod

	dbObj.GenerateVersion()
	newObj, update, err := n.IService.Apply(common.DefaultNamespace, common.AppResource, dbObj.UUID, dbObj, false)
	if err != nil {
		return nil, false, err
	}

	result := &appservice.AppResource{}
	if err = utils.UnstructuredObjectToInstanceObj(newObj, &result); err != nil {
		return nil, false, err
	}

	history.Spec.Now = result
	if _, err = n.IService.Create(common.DefaultNamespace, common.History, history); err != nil {
		return nil, false, errors.New("the history create failed")
	}
	return result, update, nil
}
