package Allconfigservice

import (
	"fmt"
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

func (a *AllConfigService) List(req *resource.Allconfigproject) {
	keys := make([]string, 0)
	for k := range req.Spec.Allconfig {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k, req.Spec.Allconfig[k])
	}
}

func (a *AllConfigService) Create(req *resource.Allconfigproject) error {
	req.GenerateVersion()
	_, err := a.IService.Create(common.DefaultNamespace, common.AllconfigProject, req)
	if err != nil {
		return err
	}
	return nil
}

func (a *AllConfigService) update(name, uuid string, value interface{}, req *resource.Allconfigproject) (core.IObject, bool, error) {
	req.GenerateVersion()
	req.Spec.Allconfig[name] = value
	return a.IService.Apply(common.DefaultNamespace, common.AllconfigProject, uuid, req)
}

func (a *AllConfigService) Delete(uuid string, singleName string) error {
	deleteElement := &resource.Allconfigproject{}
	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppProject, uuid, deleteElement); err != nil {
		return err
	}
	deleteResponse := &resource.Allconfigproject{}
	if err := utils.UnstructuredObjectToInstanceObj(deleteElement, &deleteResponse); err != nil {
		return err
	}
	delete(deleteResponse.Spec.Allconfig, singleName)
	return nil
}
