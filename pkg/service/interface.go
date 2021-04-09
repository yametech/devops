package service

import (
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/store"
)

type IService interface {
	Create(obj interface{}) error
	Save(obj interface{}) error
	Update(src interface{}, target interface{}) error
	Delete(obj interface{}) error
	List(obj interface{}) error
	Query(filter map[string]interface{}, obj interface{}) error
	Range(table string, c map[string]interface{}, f func(core.IObject) error) error
}

type BaseService struct {
	store.IStore
}

func (b *BaseService) Create(obj interface{}) error {
	err := b.Apply(obj)
	if err != nil {
		return err
	}
	return nil
}
func (b *BaseService) Save(obj interface{}) error {
	err := b.IStore.Save(obj)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseService) Update(src interface{}, target interface{}) error {
	err := b.IStore.Update(src, target)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseService) Delete(obj interface{}) error {
	err := b.Del(obj)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseService) List(obj interface{}) error {
	err := b.IStore.List(obj)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseService) Query(filter map[string]interface{}, obj interface{}) error {
	err := b.GetByFilter(filter, obj)
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
