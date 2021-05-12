package appservice

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appservice"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appservice"
	"github.com/yametech/devops/pkg/resource/workorder"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

type AppConfigService struct {
	service.IService
}

func NewAppConfigService(i service.IService) *AppConfigService {
	return &AppConfigService{IService: i}
}

func (a *AppConfigService) GetAppConfig(appid string) (core.IObject, error) {
	req := &appservice.AppConfig{}

	if err := a.IService.GetByFilter(common.DefaultNamespace, common.AppConfig, req, map[string]interface{}{
		"spec.app": appid,
	}); err != nil {
		return nil, err
	}
	return req, nil
}

func (a *AppConfigService) GetAppResources(appid string) ([]interface{}, error) {
	filter := map[string]interface{}{
		"spec.app": appid,
	}

	sort := map[string]interface{}{
		"metadata.created_time": 1,
	}

	return a.IService.ListByFilter(common.DefaultNamespace, common.AppResource, filter, sort, 0, 0)
}

func (a *AppConfigService) UpdateConfigResource(data *apiResource.NamespaceRequest) (core.IObject, bool, int, error) {



	namespaceFilter := map[string]interface{}{
		"metadata.name": data.Name,
		"spec.app":      data.App,
		"metadata.uuid": bson.M{"$ne": data.UUID},
	}

	exist := &appservice.AppResource{}
	if err := a.IService.GetByFilter(common.DefaultNamespace, common.AppResource, exist, namespaceFilter); err == nil {
		return nil, false, http.StatusBadRequest ,errors.New("this config resource is exist")
	}

	appResource := &appservice.AppResource{}
	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppResource, data.UUID, appResource); err != nil {
		log.Printf("UpdateConfigResource not found Create New one: %v\n", err)
	}

	// check the order status
	if appResource.Spec.ResourceStatus == appservice.Checking {
		return nil, false, http.StatusUnauthorized, errors.New("the workorder is checking")
	}

	// Get Cpus and Memories from cmdb
	cmdbCpus := 1000.0
	cmdbMemories := int64(10240000)

	// List all data from the same resourcePool
	appParentResource := &appservice.AppResource{}
	parentFilter := map[string]interface{}{
		"spec.app": data.ParentApp,
	}
	if err := a.IService.GetByFilter(common.DefaultNamespace, common.AppResource, appParentResource, parentFilter); err != nil {
		return nil, false, http.StatusBadRequest, errors.New("have no this namespace")
	}

	filter := map[string]interface{}{
		"spec.parent_app": appParentResource.UUID,
	}

	resources, err := a.IService.ListByFilter(common.DefaultNamespace, common.AppResource, filter, nil, 0, 0)
	if err != nil {
		log.Printf("UpdateConfigResource have no the same resourcePool use: %v", err)
	}
	appResources := make([]*appservice.AppResource, 0)
	if err = utils.UnstructuredObjectToInstanceObj(resources, &appResources); err != nil {
		return nil, false, http.StatusBadRequest, err
	}

	var (
		newCpus     float64
		newMemories int64
	)

	for _, resource := range appResources {
		newCpus += resource.Spec.Cpu
		newMemories += resource.Spec.Memory
	}

	// check Cpus and Memory
	newCpus += data.Cpu
	newMemories += data.Memory

	if appParentResource.Spec.Threshold < int((newCpus/cmdbCpus)*100) {
		return nil, false, http.StatusPaymentRequired,errors.New("The total CPU resource exceeds the limit")
	}
	if appParentResource.Spec.Threshold < int((newMemories/cmdbMemories)*100) {
		return nil, false, http.StatusPaymentRequired, errors.New("The total Memory resource exceeds the limit")
	}

	// check the other resources
	if appParentResource.Spec.Approval {
		if data.Cpu > appParentResource.Spec.Cpu {
			return nil, false, http.StatusPaymentRequired, errors.New("the Cpu resource exceeds the limit")
		}
		if data.Memory > appParentResource.Spec.Memory {
			return nil, false, http.StatusPaymentRequired, errors.New("the Memory resource exceeds the limit")
		}
		if data.Pod > appParentResource.Spec.Pod {
			return nil, false, http.StatusPaymentRequired, errors.New("the Pod resource exceeds the limit")
		}
	}

	// create history
	history := &appservice.AppConfigHistory{}
	history.Spec.Creator = ""
	history.Spec.Before = appParentResource

	// update the Resources Config
	appResource.Metadata.Name = data.Name
	appResource.Spec.Cpu = data.Cpu
	appResource.Spec.Memory = data.Memory
	appResource.Spec.Pod = data.Pod
	appResource.Spec.CpuLimit = data.CpuLimit
	appResource.Spec.MemoryLimit = data.MemoryLimit
	appResource.Spec.ParentApp = data.ParentApp
	appResource.Spec.App = data.App
	appResource.Spec.ResourceStatus = appservice.Success

	newObj, update, err := a.IService.Apply(common.DefaultNamespace, common.AppResource, appResource.UUID, appResource, true)
	if err != nil {
		return nil, false, http.StatusBadRequest, err
	}

	newAppResource := &appservice.AppResource{}
	if err = utils.UnstructuredObjectToInstanceObj(newObj, &newAppResource); err != nil {
		return nil, false, http.StatusBadRequest, err
	}

	history.Spec.Now = newAppResource
	history.Spec.App = newAppResource.Spec.App
	history.GenerateVersion()
	if _, err = a.IService.Create(common.DefaultNamespace, common.History, history); err != nil {
		return nil, false, http.StatusBadRequest, err
	}
	return newObj, update, http.StatusOK, nil
}

func (a *AppConfigService) UpdateAppConfig(data *apiResource.AppConfigRequest) (core.IObject, bool, error) {
	app := &appservice.AppProject{}
	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppProject, data.App, app); err != nil {
		return nil, false, errors.New("The app is not exist")
	}

	if app.Spec.AppType != appservice.App {
		return nil, false, errors.New("This is not an App type")
	}

	dbObj := &appservice.AppConfig{}
	err := a.IService.GetByFilter(common.DefaultNamespace, common.AppConfig, dbObj, map[string]interface{}{
		"spec.app": app.Metadata.UUID,
	})
	if err != nil {
		log.Printf("Update AppConfig Not Found Create New One: %v", err)
	}

	dbObj.Spec.Config = data.Config
	dbObj.Spec.App = data.App
	return a.IService.Apply(common.DefaultNamespace, common.AppConfig, dbObj.UUID, dbObj, true)

}

func (a *AppConfigService) DeleteResource(uuid string) error {
	return a.IService.Delete(common.DefaultNamespace, common.AppResource, uuid)
}

func (a *AppConfigService) History(appid string, page, pageSize int64) ([]interface{}, int64, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{
		"spec.app": appid,
	}

	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}

	count, err := a.IService.Count(common.DefaultNamespace, common.History, filter)
	if err != nil {
		return nil, 0, err
	}

	data, err := a.IService.ListByFilter(common.DefaultNamespace, common.History, filter, sort, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

func (a *AppConfigService) OrderToResourceCheck(obj *workorder.WorkOrder) error {

	appResource := &appservice.AppResource{}
	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppResource, obj.Spec.Relation, appResource); err != nil {
		return err
	}

	newResource := &apiResource.NamespaceRequest{}
	if err := utils.UnstructuredObjectToInstanceObj(obj.Spec.Extends, &newResource); err != nil {
		return err
	}

	appResource.Spec.ParentApp = newResource.ParentApp
	appResource.Spec.Pod = newResource.Pod
	appResource.Spec.Cpu = newResource.Cpu
	appResource.Spec.CpuLimit = newResource.CpuLimit
	appResource.Spec.Memory = newResource.Memory
	appResource.Spec.MemoryLimit = newResource.MemoryLimit
	appResource.Spec.ResourceStatus = appservice.Checking

	if _, _, err := a.IService.Apply(common.DefaultNamespace, common.AppResource, appResource.UUID, appResource, true); err != nil {
		return err
	}

	return nil
}

func (a *AppConfigService) OrderToResourceSuccess(obj *workorder.WorkOrder) error {
	appResource := &appservice.AppResource{}
	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppResource, obj.Spec.Relation, appResource); err != nil {
		return err
	}

	appResource.Spec.ResourceStatus = appservice.Success
	if _, _, err := a.IService.Apply(common.DefaultNamespace, common.AppResource, appResource.UUID, appResource, false); err != nil {
		return err
	}

	oldHistory := make([]*appservice.AppConfigHistory, 0)
	filter := map[string]interface{}{
		"spec.app": appResource.UUID,
	}
	sort := map[string]interface{}{
		"metadata.created_time": 1,
	}
	data, err := a.IService.ListByFilter(common.DefaultNamespace, common.History, filter, sort, 0, 1)
	if err != nil {
		return err
	}

	if err = utils.UnstructuredObjectToInstanceObj(data, &oldHistory); err != nil {
		return err
	}

	newHistory := &appservice.AppConfigHistory{}
	newHistory.Spec.Creator = ""
	newHistory.Spec.App = appResource.Spec.App
	newHistory.Spec.Now = appResource
	if len(oldHistory) > 0 {
		newHistory.Spec.Before = oldHistory[0].Spec.Now
	}

	newHistory.GenerateVersion()
	if _, err = a.IService.Create(common.DefaultNamespace, common.History, newHistory); err != nil {
		return err
	}

	return nil
}

func (a *AppConfigService) OrderToResourceFailed(obj *workorder.WorkOrder) error {
	appResource := &appservice.AppResource{}
	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppResource, obj.Spec.Relation, appResource); err != nil {
		return err
	}

	appResource.Spec.ResourceStatus = appservice.Failed
	if _, _, err := a.IService.Apply(common.DefaultNamespace, common.AppResource, appResource.UUID, appResource, false); err != nil {
		return err
	}

	return nil
}
