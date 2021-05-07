package artifactory

import "github.com/yametech/devops/pkg/resource/artifactory"

type RequestDeploy struct {
	DeployNamespace string                  `json:"deploy_namespace"`
	Replicas        int                     `json:"replicas"`
	AppName         string                  `json:"app_name"`                         // 只存英文名，appCode不需要，用name搜索
	DeploySpace     artifactory.DeploySpace `json:"deploy_space" bson:"deploy_space"` //部署环境

	ServicePorts  []artifactory.ServicePort  `json:"service_ports"`
	Containers    []artifactory.Container    `json:"container"`
	StorageClaims []artifactory.StorageClaim `json:"storage_claims"`
}
