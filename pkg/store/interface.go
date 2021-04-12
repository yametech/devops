package store

import (
	"fmt"
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

type IKVStore interface{}
