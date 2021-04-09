package store

type IStore interface {
	List(result interface{}) error
	//Get(table string, filter map[string]interface{}, result interface{}) error
	GetByFilter(filter map[string]interface{}, result interface{}) error
	Del(obj interface{}) error
	Apply(object interface{}) error
	Update(src, dst interface{}) error
	Save(obj interface{}) error
}
