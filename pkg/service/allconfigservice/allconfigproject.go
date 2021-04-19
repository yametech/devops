package allconfigservice

import (
	"fmt"
	resources "github.com/yametech/devops/pkg/api/resource"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	"sort"
)

type AllConfigService struct {
	service.IService
}

func NewAllConfigService(i service.IService) *AllConfigService {
	return &AllConfigService{i}
}

func (a *AllConfigService) List(search string, uuid string) (int, interface{}, interface{}, error) {
	keys := make([]string, 0)
	value := make([]interface{}, 0)
	allies := &resource.AllConfigProject{}
	restores := &resources.ConfigResponse{}
	//allies.Spec.Allconfig=make(map[string]interface{})
	if uuid != "" {
		if err := a.IService.GetByUUID(common.DefaultNamespace, common.AllconfigProject, uuid, allies); err != nil {
			return 0, nil, nil, err
		}
	}
	if err := utils.UnstructuredObjectToInstanceObj(allies, &restores); err != nil {
		return 0, nil, nil, err
	}
	if search != "" {
		if val, ok := restores.Spec.Allconfig[search]; ok {
			return 1, search, val, nil
		}
	}
	for k := range restores.Spec.Allconfig {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		value = append(value, restores.Spec.Allconfig[k])
	}

	return len(keys), keys, value, nil
}

func (a *AllConfigService) Create(req *resource.AllConfigProject) (uuid string, err error) {
	req.GenerateVersion()

	_, err = a.IService.Create(common.DefaultNamespace, common.AllconfigProject, req)
	if err != nil {
		return "", err
	}
	required := req.GetUUID()
	resourcefully := resources.ConfigResponse{}
	resourcefully.AllConfigProject.UUID = required
	return required, nil
}

func (a *AllConfigService) Update(name, uuid string, value interface{}, req *resource.AllConfigProject) (core.IObject, bool, error) {
	req.GenerateVersion()
	req.Spec.Allconfig = make(map[string]interface{})
	req.Spec.Allconfig[name] = value
	return a.IService.Apply(common.DefaultNamespace, common.AllconfigProject, uuid, req)
}

func (a *AllConfigService) Delete(uuid string, singleName string) error {
	deleteElement := &resource.AllConfigProject{}
	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AllconfigProject, uuid, deleteElement); err != nil {
		return err
	}
	deleteResponse := &resources.ConfigResponse{}

	if err := utils.UnstructuredObjectToInstanceObj(deleteElement, &deleteResponse); err != nil {
		return err
	}
	deleteResponse.Spec.Allconfig = make(map[string]interface{})
	delete(deleteResponse.Spec.Allconfig, singleName)
	for k, v := range deleteResponse.Spec.Allconfig {
		fmt.Println(k, v)
	}
	whx, res, err := a.IService.Apply(common.DefaultNamespace, common.AllconfigProject, uuid, deleteResponse)
	fmt.Println(whx, res, err)
	return nil
}
