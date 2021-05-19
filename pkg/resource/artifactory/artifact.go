package artifactory

import (
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/store"
	"github.com/yametech/devops/pkg/store/gtm"
)

type ArtifactStatus uint8

const ArtifactKind core.Kind = "artifact"

const (
	Created        ArtifactStatus = iota //未构建
	InitializeFail                       //初始化错误
	Building                             //构建中
	Built                                //构建完成
	BuiltFAIL                            //构建失败
)

type ArtifactSpec struct {
	ArtifactStatus `json:"artifact_status" bson:"artifact_status"`
	Branch         string `json:"branch" bson:"branch"`
	GitUrl         string `json:"git_url" bson:"git_url"`
	Language       string `json:"language" bson:"language"`
	Tag            string `json:"tag" bson:"tag" `
	Registry       string `json:"registry" bson:"registry"`         // 镜像库地址
	Images         string `json:"images" bson:"images"`             // CI的镜像完整路径
	ProjectPath    string `json:"project_path" bson:"project_path"` // 项目CI时WORKDIR
	ProjectFile    string `json:"project_file" bson:"project_file"` // 项目CI时Dockerfile路径
	AppName        string `json:"app_name" bson:"app_name"`         // 只存英文名，appCode不需要，用name搜索
	CreateUserId   string `json:"create_user_id" bson:"create_user_id"`
	CreateUser     string `json:"create_user" bson:"create_user"`
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

// Artifact impl Coder
func (*Artifact) Decode(op *gtm.Op) (core.IObject, error) {
	artifact := &Artifact{}
	if err := core.ObjectToResource(op.Data, artifact); err != nil {
		return nil, err
	}
	return artifact, nil
}

func init() {
	store.AddResourceCoder(string(ArtifactKind), &Artifact{})
}
