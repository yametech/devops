package appservice

import (
	"github.com/yametech/devops/pkg/core"
)

type ResourcePoolSpec struct {
	ParentApp string   `json:"parent_app" bson:"parent_app"`
	RootApp   string   `json:"root_app" bson:"root_app"`
	Desc      string   `json:"desc" bson:"desc"`
}

type ResourcePool struct {
	core.Metadata `json:"metadata"`
	Spec          ResourcePoolSpec `json:"spec"`
}

func (r *ResourcePool) Clone() core.IObject {
	result := &ResourcePool{}
	core.Clone(r, result)
	return result
}
