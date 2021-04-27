package appproject

import "github.com/yametech/devops/pkg/core"

type Resource struct {
	core.Metadata `json:"metadata"`
	Spec          ResourceSpec `json:"spec"`
}

type ResourceSpec struct {
	App          string  `json:"app" bson:"app"`
	ParentApp    string  `json:"parent_app" bson:"parent_app"`
	Cpu          float64 `json:"cpu" bson:"cpu"`
	Memory       int64   `json:"memory" bson:"memory"`
	Pod          int     `json:"pod" bson:"pod"`
	Threshold    int     `json:"threshold" bson:"threshold"`
	Cpus         float64 `json:"cpus" bson:"cpus"`
	Memories     int64   `json:"memories" bson:"memories"`
	CpuRemain    float64 `json:"cpu_remain" bson:"cpu_remain"`
	MemoryRemain int64   `json:"memory_remain" bson:"memory_remain"`
	Approval     bool    `json:"approval" bson:"approval"`
}

func (r *Resource) Clone() core.IObject {
	result := &Resource{}
	core.Clone(r, result)
	return result
}
