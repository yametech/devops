package store

import "github.com/yametech/devops/pkg/core"

type IStore interface {
	List(table string, result interface{}) error
	//Get(table string, filter map[string]interface{}, result interface{}) error
	GetByFilter(table string, filter map[string]interface{}, result interface{}) error
	Del(obj core.IObject) error
	Apply(object interface{}) error
	Update(table, uuid string, dst interface{}) error
}
