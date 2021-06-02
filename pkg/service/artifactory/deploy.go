package artifactory

import (
	"fmt"
	apiResource "github.com/yametech/devops/pkg/api/resource/artifactory"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	arResource "github.com/yametech/devops/pkg/resource/artifactory"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
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

func (a *DeployService) List(name map[string]interface{}, page, pageSize int64) ([]interface{}, int64, error) {
	offset := (page - 1) * pageSize

	for k, v := range name {
		if v == "" {
			delete(name, k)
		}
	}

	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}

	data, err := a.IService.ListByFilter(common.DefaultNamespace, common.Deploy, name, sort, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := a.IService.Count(common.DefaultNamespace, common.Deploy, name)
	if err != nil {
		return nil, 0, err
	}
	return data, count, nil

}

func (a *DeployService) GetByUUID(uuid string) (*arResource.Deploy, error) {
	dp := &arResource.Deploy{}
	err := a.IService.GetByUUID(common.DefaultNamespace, common.Deploy, uuid, dp)
	if err != nil {
		return nil, err
	}

	return dp, nil
}

func (a *DeployService) GetByAppName(appName, namespace string) (*arResource.Deploy, error) {
	dp := &arResource.Deploy{}

	filter := map[string]interface{}{
		"spec.app_name":         appName,
		"spec.artifact_status":  arResource.Deployed,
		"spec.deploy_namespace": namespace,
	}
	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}
	data, err := a.IService.ListByFilter(common.DefaultNamespace, common.Deploy, filter, sort, 0, 1)
	if err != nil {
		return nil, err
	}
	if len(data) < 1 {
		return nil, nil
	}
	err = utils.UnstructuredObjectToInstanceObj(data[0], dp)
	if err != nil {
		return nil, err
	}

	return dp, nil
}

func (a *DeployService) Create(request *apiResource.RequestDeploy) error {
	if request.Name == "" {
		request.Name = time.Now().Format("20060102-1504-05")
	}

	deploy := &arResource.Deploy{
		Metadata: core.Metadata{
			Name: request.Name,
		},
		Spec: arResource.DeploySpec{
			DeployNamespace:     request.DeployNamespace,
			DeployNamespaceUUID: request.DeployNamespaceUUID,
			Replicas:            request.Replicas,
			AppName:             request.AppName,
			DeploySpace:         request.DeploySpace,
			ServicePorts:        request.ServicePorts,
			Containers:          request.Containers,
			StorageClaims:       request.StorageClaims,
			CreateUserId:        request.UserNameID,
			CreateUser:          request.UserName,
			CreateTeam:          request.Team,
		},
	}
	for k := range deploy.Spec.Containers {
		MergeEnvVar(&deploy.Spec.Containers[k])
	}

	for k := range deploy.Spec.ServicePorts {
		deploy.Spec.ServicePorts[k].Protocol = strings.ToUpper(deploy.Spec.ServicePorts[k].Protocol)
	}

	if err := a.GetArInfo(deploy); err != nil {
		return err
	}

	deploy.GenerateVersion()
	if _, err := a.IService.Create(common.DefaultNamespace, common.Deploy, deploy); err != nil {
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

	for k, volume := range container.VolumeMount {
		configVolume := make(map[string]interface{})
		if volume.VolumeMountType == arResource.ConfigMap {
			volumeNameList := strings.Split(volume.Path, "/")
			volumeName := volumeNameList[len(volumeNameList)-1]
			volume.Key = strings.ReplaceAll(volumeName, ".", "-")
			configVolume["mountName"] = fmt.Sprintf("%s-%d", volume.Name, k)
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
		CPULimit:        Strval(container.LimitCPU),
		CPURequests:     Strval(container.RequiredCPU),
		DeployNamespace: deploy.Spec.DeployNamespace,
		DeployType:      "web", //useless key
		MemLimit:        fmt.Sprintf("%dM", container.LimitMemory),
		MemRequests:     fmt.Sprintf("%dM", container.RequiredMemory),
		Policy:          container.ImagePullPolicy,
		Replicas:        deploy.Spec.Replicas,
		ServiceImage:    container.Images,
		ServiceName:     deploy.Spec.AppName,
		StorageCapacity: "100m", //useless key
	}

	var actionName string
	switch deploy.Spec.DeploySpace {
	case arResource.SmartCity:
		actionName = common.SmartCityCD
	case arResource.Azure:
		actionName = common.AzureCD
		ReplaceRegistry(&cdInfo.ServiceImage)
	case arResource.TungChung:
		actionName = common.TungChungCD
		ReplaceRegistry(&cdInfo.ServiceImage)
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

	appName := strings.Replace(deploy.Spec.AppName, "_", "", -1)
	appName = strings.Replace(appName, "-", "", -1)
	stepName := fmt.Sprintf("%s_%s_%s", common.CD, deploy.UUID, appName)
	if !SendEchoer(stepName, actionName, sendCDInfo) {
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

func ReplaceRegistry(old *string) {
	newReg := strings.Replace(*old, "registry-d.ym", "registry-t.iauto360.cn", 1)
	*old = newReg
}

func (a *DeployService) GetArInfo(deploy *arResource.Deploy) error {

	for k, v := range deploy.Spec.Containers {
		err := a.IService.GetByUUID(common.DefaultNamespace, common.Artifactory, v.ImagesUUID, &deploy.Spec.Containers[k].ImagesInfo)
		if err != nil {
			return err
		}
	}
	return nil
}
