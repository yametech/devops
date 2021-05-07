package artifactory

import (
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/store"
	"github.com/yametech/devops/pkg/store/gtm"
)

const DeployKind core.Kind = "deploy"

type DeployStatus uint8

const (
	NotDeployed DeployStatus = iota
	Deploying
	Deployed
	DeployFail
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
	// 如果是ConfigMap 则4个都要，如果是VolumeClaim 只需要name和path
	VolumeMountType `json:"volume_mount_type" bson:"volume_mount_type"`
	Name            string `json:"name" bson:"name"`
	Path            string `json:"path" bson:"path"`
	Key             string `json:"key" bson:"key"`
	Value           string `json:"value" bson:"value"`
}

type Container struct {
	Name            string                   `json:"name" bson:"name"`
	Images          string                   `json:"images" bson:"images"`
	ImagePullPolicy string                   `json:"image_pull_policy" bson:"image_pull_policy"`
	LimitCPU        string                   `json:"limit_cpu" bson:"limit_cpu"`
	RequiredCPU     string                   `json:"required_cpu" bson:"required_cpu"`
	LimitMemory     string                   `json:"limit_memory" bson:"limit_memory"`
	RequiredMemory  string                   `json:"required_memory" bson:"required_memory"`
	Environment     []map[string]interface{} `json:"environment" bson:"environment"`
	Command         []string                 `json:"command" bson:"command"`
	Argument        []string                 `json:"argument" bson:"argument"`
	VolumeMount     []VolumeMount            `json:"volume_mounts" bson:"volume_mounts"`
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
	DeployNamespace string       `json:"deploy_namespace" bson:"deploy_namespace"`
	Replicas        int          `json:"replicas" bson:"replicas"`
	CreateUserId    string       `json:"create_user_id" bson:"create_user_id"`
	AppName         string       `json:"app_name" bson:"app_name"` // 只存英文名，appCode不需要，用name搜索
	DeployStatus    DeployStatus `json:"artifact_status" bson:"artifact_status"`
	DeploySpace     DeploySpace  `json:"deploy_space" bson:"deploy_space"` //部署环境

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
