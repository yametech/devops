package appproject

import "github.com/yametech/devops/pkg/resource/appproject"

type AppConfigRequest struct {
	App       string                 `json:"app"`
	Config    map[string]interface{} `json:"config"`
	Resources []*NameSpaceRequest    `json:"resources"`
}

type NameSpaceRequest struct {
	App       string  `json:"app" bson:"app"`
	ParentApp string  `json:"parent_app"`
	Cpu       float64 `json:"cpu" bson:"cpu"`
	Memory    int64   `json:"memory" bson:"memory"`
	Pod       int     `json:"pod" bson:"pod"`
	Threshold int     `json:"threshold" bson:"threshold"`
	Approval  bool    `json:"approval" bson:"approval"`
}

type AppConfigResponse struct {
	Config    *appproject.AppConfig  `json:"config"`
	Resources []*appproject.Resource `json:"resources"`
}
