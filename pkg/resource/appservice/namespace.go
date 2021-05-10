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

type NamespaceHistory struct {
	core.Metadata `json:"metadata"`
	Spec          NamespaceHistorySpec `json:"spec"`
}

type NamespaceHistorySpec struct {
	App     string                 `json:"app" bson:"app"`
	Creator string                 `json:"creator" bson:"creator"`
	Before  map[string]interface{} `json:"before" bson:"before"`
	Now     map[string]interface{} `json:"now" bson:"now"`
}

func (a *NamespaceHistory) Clone() core.IObject {
	result := &NamespaceHistory{}
	core.Clone(a, result)
	return result
}
