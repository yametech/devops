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
	List(table string, offset, limit int, isPreload bool, obj interface{}) (int64, error)
	Query(table string, filter map[string]interface{}, offset, limit int, isPreload bool, obj interface{}) (int64, error)
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

func (b *BaseService) List(table string, offset, limit int, isPreload bool, obj interface{}) (int64, error) {
	count, err := b.IStore.List(table, obj, offset, limit, isPreload, true)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b *BaseService) Query(table string, filter map[string]interface{}, offset, limit int, isPreload bool, obj interface{}) (int64, error) {
	count, err := b.GetByFilter(table, filter, obj, offset, limit, isPreload, true)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b *BaseService) Range(table string, c map[string]interface{}, f func(core.IObject) error) error {
	panic("implement me")
}

func NewBaseService(s store.IStore) IService {
	return &BaseService{s}
}
