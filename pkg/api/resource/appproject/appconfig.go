package appproject

import "github.com/yametech/devops/pkg/resource/appservice"

type AppConfigRequest struct {
	App    string                 `json:"app"`
	Config map[string]interface{} `json:"config"`
}

type ResourcePoolRequest struct {
	Name        string  `json:"name"`
	UUID        string  `json:"uuid"`
	App         string  `json:"app"`
	ParentApp   string  `json:"parent_app"`
	Cpu         float64 `json:"cpu"`
	Memory      int64   `json:"memory"`
	Pod         int     `json:"pod"`
	CpuLimit    float64 `json:"cpu_limit"`
	MemoryLimit int64   `json:"memory_limit"`
	Threshold   int     `json:"threshold"`
	Approval    bool    `json:"approval"`
}

type AppConfigResponse struct {
	Config    *appservice.AppConfig     `json:"config"`
	Resources []*appservice.AppResource `json:"resources"`
}
