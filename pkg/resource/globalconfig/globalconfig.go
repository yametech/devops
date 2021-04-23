package globalconfig

import "github.com/yametech/devops/pkg/core"

//type Number uint8
//
//const  (
//	ServiceSequence Number = iota
//)

type Spec struct {
	SortString []string               `json:"sort_string" bson:"sort_string"`
	Service    map[string]interface{} `json:"service" bson:"service"`
}

type GlobalConfig struct {
	core.Metadata `json:"metadata"`
	Spec          Spec `json:"spec"`
}

func (all GlobalConfig) Clone() core.IObject {
	result := &GlobalConfig{}
	core.Clone(all, result)
	return result
}
