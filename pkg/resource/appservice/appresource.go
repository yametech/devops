package appservice

import (
	"github.com/yametech/devops/pkg/core"
)

type AppResource struct {
	core.Metadata `json:"metadata"`
	Spec          AppResourceSpec `json:"spec"`
}

type ResourceStatus uint8

const (
	Success ResourceStatus = iota
	Checking
	Failed
)

type AppResourceSpec struct {
	App            string  `json:"app" bson:"app"`
	ParentApp      string  `json:"parent_app" bson:"parent_app"`
	Cpu            float64 `json:"cpu" bson:"cpu"`
	Memory         int64   `json:"memory" bson:"memory"`
	Pod            int     `json:"pod" bson:"pod"`
	CpuLimit       float64 `json:"cpu_limit" bson:"cpu_limit"`
	MemoryLimit    int64   `json:"memory_limit" bson:"memory_limit"`
	Threshold      int     `json:"threshold" bson:"threshold"`
	Approval       bool    `json:"approval" bson:"approval"`
	ResourceStatus `json:"resource_status" bson:"resource_status"`
}

func (r *AppResource) Clone() core.IObject {
	result := &AppResource{}
	core.Clone(r, result)
	return result
}
