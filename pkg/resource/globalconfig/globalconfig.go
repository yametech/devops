package globalconfig

import "github.com/yametech/devops/pkg/core"

type Spec struct {
	Service map[string]interface{} `json:"service" bson:"service"`
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
