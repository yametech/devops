package appservice

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appservice"
	"github.com/yametech/devops/pkg/service"
	"log"
)

type ResourcePoolConfigService struct {
	service.IService
}

func NewResourcePoolConfigService(i service.IService) *ResourcePoolConfigService {
	return &ResourcePoolConfigService{IService: i}
}

func (n *ResourcePoolConfigService) GetByFilter(appid string) (core.IObject, error) {
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

func (n *ResourcePoolConfigService) Update(data *apiResource.ResourcePoolRequest) (core.IObject, bool, error) {

	resourcePool := &appservice.ResourcePool{}
	if err := n.GetByUUID(common.DefaultNamespace, common.ResourcePool, data.App, resourcePool); err != nil {
		return nil, false, errors.New("The ResourcePool is not exist")
	}

	dbObj := &appservice.AppResource{}
	err := n.IService.GetByFilter(common.DefaultNamespace, common.AppResource, dbObj, map[string]interface{}{
		"spec.app": resourcePool.Metadata.UUID,
	})
	if err != nil {
		log.Printf("Update AppResource Not Found Create New One: %v\n", err)
	}


	// create history
	// Get creator
	history := &appservice.AppResourceHistory{
		Spec: appservice.AppResourceHistorySpec{
			App: dbObj.Spec.App,
			History: map[string]interface{}{
				"creator": "",
				"before":  dbObj,
			},
		},
	}

	dbObj.Spec.App = resourcePool.Metadata.UUID
	dbObj.Spec.Threshold = data.Threshold
	dbObj.Spec.Approval = data.Approval
	dbObj.Spec.Cpu = data.Cpu
	dbObj.Spec.Memory = data.Memory
	dbObj.Spec.Pod = data.Pod

	dbObj.GenerateVersion()

	result, update, err := n.IService.Apply(common.DefaultNamespace, common.AppResource, dbObj.UUID, dbObj, false)
	if err != nil {
		return nil, false, err
	}
	history.Spec.History["now"] = result
	if _, err = n.IService.Create(common.DefaultNamespace, common.History, history); err != nil {
		return nil, false, errors.New("the history create failed")
	}
	return result, update, nil
}
