package globalconfigproject

import (
	"github.com/yametech/devops/pkg/resource/globalconfig"
)

type RequestGlobalConfig struct {
	Request serverSpec `json:"request"`
	Name    string     `json:"name"`
	Kind    string     `json:"kind"`
	Content string     `json:"content"`
}

type serverSpec struct {
	Service map[string]interface{} `json:"server" bson:"server"`
}

type ConfigResponse struct {
	globalconfig.GlobalConfig
}
