package store

import "github.com/yoshiakiley/devops-zpk-server/pkg/core"

type IStore interface {
	List(db, table string, result interface{}) error
	GetByFilter(db, table string, filter map[string]interface{}, result interface{}) error
	Del(db, table string, object interface{}) error
	Apply(db, table string, object core.IObject) error
}
