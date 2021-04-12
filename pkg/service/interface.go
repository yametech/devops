package service

import (
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/store"
)

type IService interface {
	List(namespace, resource, labels string, skip, limit int64) ([]interface{}, int64, error)
	ListByFilter(namespace, resource string, filter map[string]interface{}, skip, limit int64) ([]interface{}, int64, error)
	GetByUUID(namespace, resource, uuid string, result interface{}) error
	GetByFilter(namespace, resource string, result interface{}, filter map[string]interface{}) error
	Create(namespace, resource string, object core.IObject) (core.IObject, error)
	Apply(namespace, resource, uuid string, object core.IObject) (core.IObject, bool, error)
	Delete(namespace, resource, uuid string) error
}

type BaseService struct {
	store.IKVStore
}

func NewBaseService(s store.IKVStore) IService {
	return &BaseService{s}
}
