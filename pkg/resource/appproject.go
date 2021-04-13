package resource

import "github.com/yametech/devops/pkg/core"

type AppType uint8

const (
	Service AppType = iota
	BusinessLine
	App
)

type AppSpec struct {
	AppType `json:"parent_app" bson:"parent_app"`
	Desc    string   `json:"desc" bson:"desc"`
	Owner   []string `json:"owner" bson:"owner"`
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
