package appproject

import "github.com/yametech/devops/pkg/core"

type ConfigType uint8

const (
	BaseConfig ConfigType = iota
	ArgsConfig
	FilesConfig
	ResourcesConfig
)

type AppConfig struct {
	core.Metadata `json:"metadata"`
	Spec          AppConfigSpec `json:"spec"`
}

type AppConfigSpec struct {
	App        string `json:"app" bson:"app"`
	ConfigType `json:"config_type" bson:"config_type"`
	Config     map[string]interface{} `json:"config" bson:"config"`
}

func (ap *AppConfig) Clone() core.IObject {
	result := &AppConfig{}
	core.Clone(ap, result)
	return result
}
