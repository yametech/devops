package store

import (
	"fmt"
	"github.com/yametech/devops/pkg/core"
	"gorm.io/gorm"
)

type ErrorType error

var (
	NotFound ErrorType = fmt.Errorf("notFound")
)

type IStore interface {
	List(table string, result interface{}, offset, limit int, isPreload, isDelete bool) (int64, error)
	//Get(table string, filter map[string]interface{}, result interface{}) error
	GetByFilter(table string, filter map[string]interface{}, result interface{}, offset, limit int, isPreload, isDelete bool) (int64, error)
	Del(obj interface{}) error
	Apply(object interface{}) error
	Update(src, dst interface{}) error
	Save(obj interface{}) error
	Dao(table string) *gorm.DB
}

type IKVStore interface {
	List(namespace, resource, labels string, sort map[string]interface{}, skip, limit int64) ([]interface{}, error)
	ListByFilter(namespace, resource string, filter, sort map[string]interface{}, skip, limit int64) ([]interface{}, error)
	GetByUUID(namespace, resource, uuid string, result interface{}) error
	GetByFilter(namespace, resource string, result interface{}, filter map[string]interface{}) error
	Create(namespace, resource string, object core.IObject) (core.IObject, error)
	Apply(namespace, resource, name string, object core.IObject) (core.IObject, bool, error)
	Delete(namespace, resource, uuid string) error
	Count(namespace, resource string, filter map[string]interface{}) (int64, error)
}
