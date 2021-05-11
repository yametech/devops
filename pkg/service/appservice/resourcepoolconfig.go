package appservice

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appservice"
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

	dbObj.Spec.App = namespace.Metadata.UUID
	dbObj.Spec.Threshold = data.Threshold
	dbObj.Spec.Approval = data.Approval
	dbObj.Spec.Cpu = data.Cpu
	dbObj.Spec.Memory = data.Memory
	dbObj.Spec.Pod = data.Pod

	dbObj.GenerateVersion()
	return n.IService.Apply(common.DefaultNamespace, common.AppResource, dbObj.UUID, dbObj, false)
}

func (n *ResourcePoolConfigService) GetNamespaceResourceRemain(appid string) (float64, int64, error) {
	// get namespace all resource from cmdb
	cmdbCpus := 10000.0
	cmdbMemories := int64(10000000)

	filter := map[string]interface{}{
		"spec.parent_app": appid,
	}

	data, err := n.IService.ListByFilter(common.DefaultNamespace, common.AppResource, filter, nil, 0, 0)
	if err != nil {
		return cmdbCpus, cmdbMemories, err
	}

	children := make([]*appservice.AppResource, 0)
	if err = utils.UnstructuredObjectToInstanceObj(data, &children); err != nil {
		return cmdbCpus, cmdbMemories, err
	}

	var (
		useCpu      float64
		useMemories int64
	)
	for _, child := range children {
		useCpu += child.Spec.Cpu
		useMemories += child.Spec.Memory
	}

	return cmdbCpus - useCpu, cmdbMemories - useMemories, nil
}

func (n *ResourcePoolConfigService) GetNamespaceResource(appid string) (float64, int64, float64, float64, error) {

	// get namespace all resource from cmdb
	cmdbCpus := 10000.0
	cmdbMemories := int64(10000000)
	moneyMonth := 12000.23
	moneyYear := 123123.22

	return cmdbCpus, cmdbMemories, moneyMonth, moneyYear, nil
}
