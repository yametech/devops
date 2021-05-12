package base

import "github.com/yametech/devops/pkg/core"

type Module struct {
	core.Metadata `json:"metadata"`
	Spec          ModuleSpec `json:"spec"`
}

type ModuleSpec struct {
	Parent  string                 `json:"parent" bson:"parent"`
	Extends map[string]interface{} `json:"extends" bson:"extends"`
}

func (m *Module) Clone() core.IObject {
	result := &Module{}
	core.Clone(m, result)
	return result
}

type PrivateModule struct {
	core.Metadata `json:"metadata"`
	Spec          PrivateModuleSpec `json:"spec"`
}

type PrivateModuleSpec struct {
	Modules []string `json:"modules" bson:"modules"`
	User    string   `json:"user" bson:"user"`
}

func (p *PrivateModule) Clone() core.IObject {
	result := &PrivateModule{}
	core.Clone(p, result)
	return result
}
