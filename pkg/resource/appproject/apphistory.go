package appproject

import "github.com/yametech/devops/pkg/core"

type ConfigHistory struct {
	core.Metadata `json:"metadata"`
	Spec          HistorySpec `json:"spec"`
}

type HistorySpec struct {
	App     string                 `json:"app" bson:"app"`
	History map[string]interface{} `json:"history" bson:"history"`
}

func (c *ConfigHistory) Clone() core.IObject {
	result := &ConfigHistory{}
	core.Clone(c, result)
	return result
}
