package base

import (
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
)

type CollectionModuleService struct {
	service.IService
}

func NewCollectionModuleService(i service.IService) *CollectionModuleService {
	return &CollectionModuleService{IService: i}
}

func (c *CollectionModuleService) AddCollectionModule(uuid string) ([]interface{}, error) {
	// Check the uuid
	module := &base.Module{}
	if err := c.IService.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, module); err != nil {
		return nil, errors.New("The module uuid is not exist")
	}

	// Get the user
	user := ""

	dbObj := &base.PrivateModule{}
	if err := c.IService.GetByFilter(common.DefaultNamespace, common.CollectionModule, dbObj, map[string]interface{}{
		"spec.user": user,
	}); err != nil {
		dbObj.Spec.Modules = make([]string, 0)
		dbObj.Spec.User = user
	}

	dbObj.Spec.Modules = append(dbObj.Spec.Modules, uuid)
	if _, _, err := c.IService.Apply(common.DefaultNamespace, common.CollectionModule, dbObj.UUID, dbObj, true); err != nil {
		return nil, err
	}

	return c.ListCollectionModule()
}

func (c *CollectionModuleService) ListCollectionModule() ([]interface{}, error) {
	// Get the user
	user := ""

	collection := &base.PrivateModule{}
	response := make([]interface{}, 0)
	if err := c.IService.GetByFilter(common.DefaultNamespace, common.CollectionModule, collection, map[string]interface{}{
		"spec.user": user,
	}); err != nil {
		return response, nil
	}


	for _, uuid := range collection.Spec.Modules{
		data := &base.Module{}
		if err := c.IService.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, data); err != nil {
			return nil, err
		}
		response = append(response, data)
	}

	return response, nil
}

func (c *CollectionModuleService) DeleteCollectionModule(uuid string) ([]interface{}, error) {
	// Get the user
	user := ""

	collection := &base.PrivateModule{}
	if err := c.IService.GetByFilter(common.DefaultNamespace, common.CollectionModule, collection, map[string]interface{}{
		"spec.user": user,
	}); err != nil {
		return c.ListCollectionModule()
	}

	for i := len(collection.Spec.Modules) - 1; i >= 0; i-- {
		if collection.Spec.Modules[i] == uuid {
			collection.Spec.Modules = append(collection.Spec.Modules[:i], collection.Spec.Modules[i+1:]...)
		}
	}

	if _, _, err := c.IService.Apply(common.DefaultNamespace, common.CollectionModule, collection.UUID, collection, true); err != nil {
		return nil, err
	}

	return c.ListCollectionModule()
}
