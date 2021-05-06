package appservice

import (
	"github.com/yametech/devops/pkg/core"
)

type NamespaceSpec struct {
	ParentApp string `json:"parent_app" bson:"parent_app"`
	Desc      string `json:"desc" bson:"desc"`
}

type Namespace struct {
	core.Metadata `json:"metadata"`
	Spec          NamespaceSpec `json:"spec"`
}

func (r *Namespace) Clone() core.IObject {
	result := &Namespace{}
	core.Clone(r, result)
	return result
}
