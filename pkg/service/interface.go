package service

import (
	"github.com/yametech/devops-zpk-server/pkg/core"
	"github.com/yametech/devops-zpk-server/pkg/store"
)

type IService interface {
	Create(db, table string, obj core.IObject) error
	Update(db, table string, src core.IObject, target core.IObject) error
	Delete(db, table string, uuid string) error
	Query(db, table string, c map[string]interface{}) ([]core.IObject, error)
	QueryOne(db, table string, c map[string]interface{}) (core.IObject, error)
	Range(db, table string, c map[string]interface{}, f func(core.IObject) error) error
}

type BaseService struct {
	store.IStore
}

func (b *BaseService) Create(db, table string, obj core.IObject) error {
	err := b.Apply(db, table, obj)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseService) Update(db, table string, src core.IObject, target core.IObject) error {
	panic("implement me")
}

func (b *BaseService) Delete(db, table string, uuid string) error {
	panic("implement me")
}

func (b *BaseService) Query(db, table string, c map[string]interface{}) ([]core.IObject, error) {
	panic("implement me")
}

func (b *BaseService) QueryOne(db, table string, c map[string]interface{}) (core.IObject, error) {
	panic("implement me")
}

func (b *BaseService) Range(db, table string, c map[string]interface{}, f func(core.IObject) error) error {
	panic("implement me")
}

func NewBaseService(s store.IStore) IService {
	return &BaseService{s}
}
