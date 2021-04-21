package artifactory

import (
	"github.com/yametech/devops/pkg/core"
)

type ArtifactStatus uint8

const ArtifactKind core.Kind = "artifact"

const (
	Created ArtifactStatus = iota
	Building
	Built
	BuiltFAIL
)

type ArtifactSpec struct {
	ArtifactStatus `json:"artifact_status" bson:"artifact_status"`
	Branch         string `json:"branch" bson:"branch"`
	GitUrl         string `json:"git_url" bson:"git_url"`
	Language       string `json:"language" bson:"language"`
	Tag            string `json:"tag" bson:"tag" `
	Images         string `json:"images" bson:"images"`
	AppConfig      string `json:"app_config" bson:"app_config"` // gitUrl和Language存在全局配置中，不需要保存在这里
	AppName        string `json:"app_name" bson:"app_name"`     // 只存英文名，appCode不需要，用name搜索
	CreateUserId   string `json:"create_user_id" bson:"create_user_id"`
	Remarks        string `json:"remarks" bson:"remarks"`
}

type Artifact struct {
	core.Metadata `json:"metadata"`
	Spec          ArtifactSpec `json:"spec"`
}

func (ar *Artifact) Clone() core.IObject {
	result := &Artifact{}
	core.Clone(ar, result)
	return result
}
