package resource

import "github.com/yametech/devops/pkg/core"

type Allspec struct {
	Allconfig map[string]interface{}
}

type Allconfigproject struct {
	core.Metadata `json:"metadata"`
	Spec          Allspec `json:"spec"`
}

func (all Allconfigproject) Clone() core.IObject {
	result := &Allconfigproject{}
	core.Clone(all, result)
	return result
}
