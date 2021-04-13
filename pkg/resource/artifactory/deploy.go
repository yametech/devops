package artifactory

import (
	"github.com/yametech/devops/pkg/core"
)

type DeployStatus uint8

const (
	NotDeployed DeployStatus = iota
	Deploying
	Deployed
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
	Images          string                   `json:"images" bson:"images"`
	ImagePullPolicy string                   `json:"image_pull_policy"`
	LimitCPU        string                   `json:"limit_cpu" bson:"limit_cpu"`
	RequiredCPU     string                   `json:"required_cpu" bson:"required_cpu"`
	LimitMemory     string                   `json:"limit_memory" bson:"limit_memory"`
	RequiredMemory  string                   `json:"required_memory" bson:"required_memory"`
	Environment     []map[string]interface{} `json:"environment" bson:"environment"`
	Command         []string                 `json:"command" bson:"command"`
	Argument        []string                 `json:"argument" bson:"argument"`
	VolumeMounts    []VolumeMount            `json:"volume_mounts" bson:"volume_mounts"`
}

type ServicePort struct {
	Name       string `json:"name" bson:"name"`
	Protocol   string `json:"protocol" bson:"protocol"`
	Port       string `json:"port" bson:"port"`
	TargetPort string `json:"target_port" bson:"target_port"`
}

type DeploySpec struct {
	DeployNamespace string `json:"deploy_namespace" bson:"deploy_namespace"`
	DeployType      string `json:"deploy_type" bson:"deploy_type"`
	StorageCapacity string `json:"storage_capacity" bson:"storage_capacity"`
	Replicas        string `json:"replicas" bson:"replicas"`
	Policy          string `json:"policy" bson:"policy"`
	CreateUserId    string `json:"create_user_id" bson:"create_user_id"`
	AppName         string `json:"app_name" bson:"app_name"` // 只存英文名，appCode不需要，用name搜索

	ServicePorts []ServicePort `json:"service_ports"`
	DeployStatus DeployStatus  `json:"artifact_status" bson:"artifact_status"`
	Container    Container     `json:"container"`
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
