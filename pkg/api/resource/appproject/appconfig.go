package appproject

import "github.com/yametech/devops/pkg/resource/appproject"

type AppConfigRequest struct {
	App       string                       `json:"app"`
	Config    map[string]interface{}       `json:"config"`
	Resources []*NameSpaceRequest          `json:"resources"`
	Totals    map[string]*NameSpaceRequest `json:"totals"`
}

type NameSpaceRequest struct {
	Name        string `json:"name"`
	UUID      string  `json:"uuid"`
	App       string  `json:"app"`
	ParentApp string  `json:"parent_app"`
	Cpu       float64 `json:"cpu"`
	Memory    int64   `json:"memory"`
	Pod       int     `json:"pod"`
	Threshold int     `json:"threshold"`
	Approval  bool    `json:"approval"`
}

type AppConfigResponse struct {
	Config    *appproject.AppConfig  `json:"config"`
	Resources []*appproject.Resource `json:"resources"`
}
