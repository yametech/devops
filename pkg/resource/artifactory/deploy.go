package artifactory

import (
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/store"
	"github.com/yametech/devops/pkg/store/gtm"
)

const DeployKind core.Kind = "deploy"

type DeployStatus uint8

const (
	NotDeployed DeployStatus = iota //未部署
	Deploying                       //部署中
	Deployed                        //部署成功
	DeployFail                      //部署失败
)

type DeploySpace uint8

const (
	SmartCity DeploySpace = iota //智慧城
	Azure                        //微软云(测试环境)
	TungChung                    //东涌(生产环境)
)

type VolumeMountType uint8

const (
	ConfigMap VolumeMountType = iota
	VolumeClaim
)

type VolumeMount struct {
	// 如果是ConfigMap 则4个都要，如果是VolumeClaim 只需要path和value
	VolumeMountType `json:"volume_mount_type" bson:"volume_mount_type"`
	Name            string `json:"name" bson:"name"`
	Path            string `json:"path" bson:"path"`
	Key             string `json:"key" bson:"key"`
	Value           string `json:"value" bson:"value"`
}

type Container struct {
	Name            string                   `json:"name" bson:"name"`
	AppUUID         string                   `json:"app_uuid" bson:"app_uuid"`
	Images          string                   `json:"images" bson:"images"`           //镜像路径
	ImagesUUID      string                   `json:"images_uuid" bson:"images_uuid"` //镜像uuid
	ImagePullPolicy string                   `json:"image_pull_policy" bson:"image_pull_policy"`
	LimitCPU        float64                  `json:"limit_cpu" bson:"limit_cpu"`
	RequiredCPU     float64                  `json:"required_cpu" bson:"required_cpu"`
	LimitMemory     int                      `json:"limit_memory" bson:"limit_memory"`
	RequiredMemory  int                      `json:"required_memory" bson:"required_memory"`
	Environment     []map[string]interface{} `json:"environment" bson:"environment"`
	Command         []string                 `json:"command" bson:"command"`
	Argument        []string                 `json:"argument" bson:"argument"`
	VolumeMount     []VolumeMount            `json:"volume_mounts" bson:"volume_mounts"`
	ImagesInfo      Artifact                 `json:"images_info,omitempty" bson:"images_info,omitempty"`
}

type ServicePort struct {
	Name       string `json:"name" bson:"name"`
	Protocol   string `json:"protocol" bson:"protocol"`
	Port       int    `json:"port" bson:"port"`
	TargetPort int    `json:"targetPort" bson:"targetPort"`
}

type StorageClaim struct {
	Name        string `json:"name" bson:"name"`
	StorageSize string `json:"storage_size" bson:"storage_size"`
}

type DeploySpec struct {
	CreateTeam          string       `json:"create_team" bson:"create_team"`
	CreateUser          string       `json:"create_user" bson:"create_user"`
	CreateUserId        string       `json:"create_user_id" bson:"create_user_id"`
	DeployNamespace     string       `json:"deploy_namespace" bson:"deploy_namespace"`
	DeployNamespaceUUID string       `json:"deploy_namespace_uuid" bson:"deploy_namespace_uuid"` //命名空间uuid，前端回填用
	Replicas            int          `json:"replicas" bson:"replicas"`
	AppName             string       `json:"app_name" bson:"app_name"` // 只存英文名，appCode不需要，用name搜索
	DeployStatus        DeployStatus `json:"artifact_status" bson:"artifact_status"`

	DeploySpace   DeploySpace    `json:"deploy_space" bson:"deploy_space"` //部署环境
	ServicePorts  []ServicePort  `json:"service_ports" bson:"service_ports"`
	Containers    []Container    `json:"containers" bson:"containers"`         //one pod may have multiple container
	StorageClaims []StorageClaim `json:"storage_claims" bson:"storage_claims"` //pod may mount multiple storage
}

type Deploy struct {
	core.Metadata `json:"metadata"`
	Spec          DeploySpec `json:"spec"`
}

func (d *Deploy) Clone() core.IObject {
	result := &Deploy{}
	core.Clone(d, result)
	return result
}

// Deploy impl Coder
func (*Deploy) Decode(op *gtm.Op) (core.IObject, error) {
	deploy := &Deploy{}
	if err := core.ObjectToResource(op.Data, deploy); err != nil {
		return nil, err
	}
	return deploy, nil
}

func init() {
	store.AddResourceCoder(string(DeployKind), &Deploy{})
}
