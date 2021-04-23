package globalconfigproject

import (
	"github.com/yametech/devops/pkg/resource/globalconfig"
)

const RequestGlobalConfigUUID = "12345678"

type RequestGlobalConfig struct {
	SortString []string               `json:"sort_string" bson:"sort_string"`
	Service    map[string]interface{} `json:"service" bson:"service"`
}

type ConfigResponse struct {
	globalconfig.GlobalConfig
}
