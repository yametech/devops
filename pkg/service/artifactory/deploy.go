package artifactory

import (
	"fmt"
	apiResource "github.com/yametech/devops/pkg/api/resource/artifactory"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	arResource "github.com/yametech/devops/pkg/resource/artifactory"
	"github.com/yametech/devops/pkg/service"
	"strings"
	"time"
)

type DeployService struct {
	service.IService
}

func NewDeployService(i service.IService) *DeployService {
	return &DeployService{i}
}

func (a *DeployService) Watch(version string) (chan core.IObject, chan struct{}) {
	objectChan := make(chan core.IObject, 32)
	closed := make(chan struct{})
	a.IService.Watch(common.DefaultNamespace, common.Deploy, string(arResource.DeployKind), version, objectChan, closed)
	return objectChan, closed
}

func (a *DeployService) List(name string, page, pageSize int64) ([]interface{}, int64, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	if name != "" {
		filter["spec.app_name"] = name
	}
	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}

	data, err := a.IService.ListByFilter(common.DefaultNamespace, common.Deploy, filter, sort, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := a.IService.Count(common.DefaultNamespace, common.Deploy, filter)
	if err != nil {
		return nil, 0, err
	}
	return data, count, nil

}

func (a *DeployService) Create(request *apiResource.RequestDeploy) error {
	deploy := &arResource.Deploy{
		Metadata: core.Metadata{
			Name: fmt.Sprintf("%s-%d", request.AppName, time.Now().Unix()),
		},
		Spec: arResource.DeploySpec{
			DeployNamespace: request.DeployNamespace,
			Replicas:        request.Replicas,
			AppName:         request.AppName,
			ServicePorts:    request.ServicePorts,
			Containers:      request.Containers,
			StorageClaims:   request.StorageClaims,
		},
	}

	deploy.GenerateVersion()
	_, err := a.IService.Create(common.DefaultNamespace, common.Deploy, deploy)
	if err != nil {
		return err
	}
	go a.sendCD(deploy)
	return nil
}

func (a *DeployService) sendCD(deploy *arResource.Deploy) {
	// 目前只支持一个container的部署方式
	//TODO 改版deploy, 现在的版本跟dxp一样丑
	container := deploy.Spec.Containers[0]
	configVolumes := make([]map[string]interface{}, 0)

	for _, volume := range container.VolumeMount {
		configVolume := make(map[string]interface{})
		if volume.VolumeMountType == arResource.ConfigMap {
			volumeNameList := strings.Split(volume.Path, "/")
			volumeName := volumeNameList[len(volumeNameList)-1]

			configVolume["mountName"] = volume.Name
			configVolume["mountPath"] = volume.Path
			configVolume["kind"] = "configmap"
			configVolume["cmItems"] = []map[string]interface{}{
				{
					"volumeName": volume.Key,
					"volumePath": volumeName,
					"volumeData": volume.Value,
				},
			}
		} else {
			configVolume["mountPath"] = volume.Path
			configVolume["kind"] = "storage"
			configVolume["mountName"] = volume.Value
		}

		configVolumes = append(configVolumes, configVolume)
	}

	artifactInfo := map[string]interface{}{
		"servicePorts":  deploy.Spec.ServicePorts,
		"command":       container.Command,
		"arguments":     container.Argument,
		"environments":  container.Environment,
		"configVolumes": configVolumes,
	}

	cdInfo := &arResource.ArtifactCDInfo{
		ArtifactInfo:    artifactInfo,
		CPULimit:        container.LimitCPU,
		CPURequests:     container.RequiredCPU,
		DeployNamespace: deploy.Spec.DeployNamespace,
		DeployType:      "web", //useless key
		MemLimit:        container.LimitMemory,
		MemRequests:     container.RequiredMemory,
		Policy:          container.ImagePullPolicy,
		Replicas:        deploy.Spec.Replicas,
		ServiceImage:    container.Images,
		ServiceName:     deploy.Spec.AppName,
		StorageCapacity: "100m", //useless key
	}
	sendCDInfo, err := core.ToMap(cdInfo)
	if err != nil {
		deploy.Spec.DeployStatus = arResource.DeployFail
		_, _, err = a.IService.Apply(common.DefaultNamespace, common.Deploy, deploy.UUID, deploy, false)
		if err != nil {
			fmt.Printf("sendcd initialize save error %s", err)
		}
		return
	}
	if !SendEchoer(deploy.UUID, common.EchoerCD, sendCDInfo) {
		deploy.Spec.DeployStatus = arResource.DeployFail
		_, _, err = a.IService.Apply(common.DefaultNamespace, common.Deploy, deploy.UUID, deploy, false)
		if err != nil {
			fmt.Printf("sendcd sendEchoer save error %s", err)
		}
		return
	}
	deploy.Spec.DeployStatus = arResource.Deploying

	_, _, err = a.IService.Apply(common.DefaultNamespace, common.Deploy, deploy.UUID, deploy, false)
	if err != nil {
		fmt.Printf("sendcd sendEchoer success save error %s", err)
	}

}
