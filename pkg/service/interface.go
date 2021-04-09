package service

import (
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/store"
)

type IService interface {
	Create(table string, obj interface{}) error
	Update(table, uuid string, target interface{}) error
	Delete(table string, obj core.IObject) error
	Query(table string, c map[string]interface{}) (interface{}, error)
	QueryOne(table string, c map[string]interface{}, obj core.IObject) error
	Range(table string, c map[string]interface{}, f func(core.IObject) error) error
}

type BaseService struct {
	store.IStore
}

func (b *BaseService) Create(table string, obj interface{}) error {
	err := b.Apply(obj)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseService) Update(table, uuid string, target interface{}) error {
	err := b.IStore.Update(table, uuid, target)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseService) Delete(table string, obj core.IObject) error {
	obj.Delete()
	err := b.Update(table, obj.GetUUID(), obj)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseService) Query(table string, c map[string]interface{}) (interface{}, error) {
	var results []map[string]interface{}
	if c != nil {
		err := b.GetByFilter(table, c, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	}
	err := b.List(table, &results)
	if err != nil {
		return nil, err
	}
	return results, nil

}

func (b *BaseService) QueryOne(table string, c map[string]interface{}, obj core.IObject) error {
	if c != nil {
		err := b.GetByFilter(table, c, obj)
		if err != nil {
			return err
		}
		return nil
	}
	err := b.GetByFilter(table, nil, obj)
	if err != nil {
		return err
	}
	return nil

}

func (b *BaseService) Range(table string, c map[string]interface{}, f func(core.IObject) error) error {
	panic("implement me")
}

func NewBaseService(s store.IStore) IService {
	return &BaseService{s}
}
