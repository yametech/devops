package appproject

import "github.com/yametech/devops/pkg/resource/appproject"

type AppConfigRequest struct {
	App     string                 `json:"app"`
	Config  map[string]interface{} `json:"config"`
	AppType appproject.AppType     `json:"app_type"`
}
