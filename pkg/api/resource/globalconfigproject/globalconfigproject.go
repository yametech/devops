package globalconfigproject

import (
	"github.com/yametech/devops/pkg/resource/globalconfig"
)

type RequestGlobalConfig struct {
	Service map[string]interface{} `json:"service" bson:"service"`
}

type ConfigResponse struct {
	globalconfig.GlobalConfig
}
