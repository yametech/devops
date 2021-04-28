package appservice

import (
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/workorder"
)

type AppResource struct {
	core.Metadata `json:"metadata"`
	Spec          AppResourceSpec `json:"spec"`
}

type AppResourceSpec struct {
	App                   string  `json:"app" bson:"app"`
	ParentApp             string  `json:"parent_app" bson:"parent_app"`
	Cpu                   float64 `json:"cpu" bson:"cpu"`
	Memory                int64   `json:"memory" bson:"memory"`
	Pod                   int     `json:"pod" bson:"pod"`
	CpuLimit              float64 `json:"cpu_limit" bson:"cpu_limit"`
	MemoryLimit           int64   `json:"memory_limit" bson:"memory_limit"`
	Threshold             int     `json:"threshold" bson:"threshold"`
	Approval              bool    `json:"approval" bson:"approval"`
	workorder.OrderStatus `json:"order_status" bson:"order_status"`
}

func (r *AppResource) Clone() core.IObject {
	result := &AppResource{}
	core.Clone(r, result)
	return result
}

type AppResourceHistory struct {
	core.Metadata `json:"metadata"`
	Spec          AppResourceHistorySpec `json:"spec"`
}

type AppResourceHistorySpec struct {
	App     string                 `json:"app" bson:"app"`
	History map[string]interface{} `json:"history" bson:"history"`
}

func (a *AppResourceHistory) Clone() core.IObject {
	result := &AppResourceHistory{}
	core.Clone(a, result)
	return result
}
