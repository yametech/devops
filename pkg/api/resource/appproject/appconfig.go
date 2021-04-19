package appproject

import "github.com/yametech/devops/pkg/resource/appproject"

type AppConfigRequest struct {
	App                   string `json:"app"`
	appproject.ConfigType `json:"config_type"`
	Config                map[string]interface{} `json:"config"`
}
