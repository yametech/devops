package resource

import "github.com/yametech/devops/pkg/core"

type AppType uint8

const (
	BusinessLine AppType = iota
	Service
	App
)

type AppProjectResponse struct {
	AppProject
	Children []*AppProjectResponse `json:"children"`
}

type AppSpec struct {
	AppType   `json:"app_type" bson:"app_type"`
	ParentApp string   `json:"parent_app" bson:"parent_app"`
	Desc      string   `json:"desc" bson:"desc"`
	Owner     []string `json:"owner" bson:"owner"`
}

type AppProject struct {
	core.Metadata `json:"metadata"`
	Spec          AppSpec `json:"spec"`
}

func (ap *AppProject) Clone() core.IObject {
	result := &AppProject{}
	core.Clone(ap, result)
	return result
}
