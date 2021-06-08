package artifactory

import "github.com/yametech/devops/pkg/resource/artifactory"

type RequestDeploy struct {
	Team                string `json:"team"`
	Name                string `json:"name"`
	UserName            string `json:"user_name"`
	UserNameID          string `json:"user_name_id"`
	DeployNamespaceUUID string `json:"deploy_namespace_uuid"`

	DeployNamespace string `json:"deploy_namespace"`
	Replicas        int    `json:"replicas"`
	AppName         string `json:"app_name"` // 只存英文名，appCode不需要，用name搜索

	DeploySpace   artifactory.DeploySpace    `json:"deploy_space" bson:"deploy_space"` //部署环境
	ServicePorts  []artifactory.ServicePort  `json:"service_ports"`
	Containers    []artifactory.Container    `json:"container"`
	StorageClaims []artifactory.StorageClaim `json:"storage_claims"`
}
