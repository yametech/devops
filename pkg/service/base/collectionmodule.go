package base

import (
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
)

type CollectionModuleService struct {
	service.IService
}

func NewCollectionModuleService(i service.IService) *CollectionModuleService {
	return &CollectionModuleService{IService: i}
}

func (c *CollectionModuleService) AddCollectionModule(uuid string, user string) (core.IObject, bool, error) {
	// Check the uuid
	module := &base.Module{}
	if err := c.IService.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, module); err != nil {
		return nil, false, errors.New("The module uuid is not exist")
	}

	dbObj := &base.PrivateModule{}
	if err := c.IService.GetByFilter(common.DefaultNamespace, common.CollectionModule, dbObj, map[string]interface{}{
		"spec.user": user,
	}); err != nil {
		dbObj.Spec.Modules = make([]string, 0)
		dbObj.Spec.User = user
	}

	for _, dbModule := range dbObj.Spec.Modules {
		if dbModule == uuid{
			return nil, false, errors.New("The module uuid is collected by this user")
		}
	}

	dbObj.Spec.Modules = append(dbObj.Spec.Modules, uuid)
	return c.IService.Apply(common.DefaultNamespace, common.CollectionModule, dbObj.UUID, dbObj, true)
}

func (c *CollectionModuleService) ListCollectionModule(user string) ([]interface{}, error) {

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

func (c *CollectionModuleService) DeleteCollectionModule(uuid string, user string) (core.IObject, bool, error) {

	collection := &base.PrivateModule{}
	if err := c.IService.GetByFilter(common.DefaultNamespace, common.CollectionModule, collection, map[string]interface{}{
		"spec.user": user,
	}); err != nil {
		return nil, false, err
	}

	for i := len(collection.Spec.Modules) - 1; i >= 0; i-- {
		if collection.Spec.Modules[i] == uuid {
			collection.Spec.Modules = append(collection.Spec.Modules[:i], collection.Spec.Modules[i+1:]...)
		}
	}

	return c.IService.Apply(common.DefaultNamespace, common.CollectionModule, collection.UUID, collection, true)
}
