package resource

import "github.com/yametech/devops/pkg/core"

type Allspec struct {
	Allconfig map[string]interface{}
}

type AllConfigProject struct {
	core.Metadata `json:"metadata"`
	Spec          Allspec `json:"spec"`
}

func (all AllConfigProject) Clone() core.IObject {
	result := &AllConfigProject{}
	core.Clone(all, result)
	return result
}
