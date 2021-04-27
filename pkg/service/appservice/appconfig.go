package appservice

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/resource/appproject"
	"github.com/yametech/devops/pkg/resource/workorder"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	apiWorkorder "github.com/yametech/devops/pkg/api/resource/workorder"
)

type AppConfigService struct {
	service.IService
}

func NewAppConfigService(i service.IService) *AppConfigService {
	return &AppConfigService{IService: i}
}

func (a *AppConfigService) GetByFilter(appid string) (*apiResource.AppConfigResponse, error) {
	req := &appproject.AppConfig{
		Spec: appproject.AppConfigSpec{
			App: appid,
		},
	}

	if err := a.IService.GetByFilter(common.DefaultNamespace, common.AppConfig, req, map[string]interface{}{
		"spec.app": req.Spec.App,
	}); err != nil {
		return nil, err
	}

	resource := make([]*appproject.Resource, 0)
	filter := map[string]interface{}{
		"spec.app": appid,
	}
	sort := map[string]interface{}{
		"metadata.create_time": 1,
	}

	data, err := a.IService.ListByFilter(common.DefaultNamespace, common.Resource, filter, sort, 0, 0)
	if err != nil {
		return nil, err
	}

	if err = utils.UnstructuredObjectToInstanceObj(data, &resource); err != nil {
		return nil, err
	}

	response := apiResource.AppConfigResponse{
		Config:    req,
		Resources: resource,
	}

	return &response, nil
}

func (a *AppConfigService) Update(data *apiResource.AppConfigRequest) (*apiResource.AppConfigResponse, bool, error) {

	// check the workorder status
	url := fmt.Sprintf("%s?relation=%s&order_type=%d", common.WorkOrderStatus, data.App, workorder.Resources)
	resp, err := utils.Request("GET", url, nil, nil)
	if err != nil {
		return nil, false, errors.New("Can not Get the worker order status")
	}

	status := &apiWorkorder.WorkOrderStatusResponse{}
	if err = json.Unmarshal(resp, &status); err != nil {
		return nil, false, err
	}

	if status.Data == workorder.Checking {
		return nil, false, errors.New("the worker order is checking, can not submit the data")
	}

	// check Cpus and Memories
	checkMap := make(map[string]*appproject.Resource)
	for parentId, total := range data.Totals {
		parent := &appproject.Resource{}
		if err := a.IService.GetByUUID(common.DefaultNamespace, common.Resource, parentId, parent); err != nil {
			return nil, false, errors.New("The resource is not exist")
		}
		checkMap[parentId] = parent
		limitCpu := int(((parent.Spec.Cpus - (parent.Spec.CpuRemain - total.Cpu)) / parent.Spec.Cpus) * 100)
		if limitCpu > parent.Spec.Threshold {
			err := fmt.Sprintf("CPU total resources exceed %d", parent.Spec.Threshold)
			return nil, false, errors.New(err)
		}
		limitMemory := int(((parent.Spec.Memories - (parent.Spec.MemoryRemain - total.Memory)) / parent.Spec.Memory) * 100)
		if limitMemory > parent.Spec.Threshold {
			err := fmt.Sprintf("Memory total resources exceed %d", parent.Spec.Threshold)
			return nil, false, errors.New(err)
		}
	}

	// check the other resources
	for _, resource := range data.Resources {
		parent := checkMap[resource.ParentApp]
		if !parent.Spec.Approval {
			continue
		}
		if resource.Cpu > parent.Spec.Cpu {
			return nil, false, errors.New("The CPU resource exceeds limit")
		}
		if resource.Memory > parent.Spec.Memory {
			return nil, false, errors.New("The Memory resource exceeds limit")
		}
		if resource.Pod > parent.Spec.Pod {
			return nil, false, errors.New("The Pod resource exceeds limit")
		}
	}

	// update the Resources Config
	updateResource := make([]*appproject.Resource, 0)
	for _, resource := range data.Resources {
		parent := checkMap[resource.ParentApp]
		res := &appproject.Resource{}
		a.IService.GetByUUID(common.DefaultNamespace, common.Resource, resource.UUID, res)

		originParent := &appproject.Resource{}
		if err := a.IService.GetByUUID(common.DefaultNamespace, common.Resource, res.Spec.ParentApp, originParent); err == nil {
			originParent.Spec.CpuRemain += res.Spec.Cpu
			originParent.Spec.MemoryRemain += res.Spec.Memory

			originParent.GenerateVersion()
			if _, _, err := a.IService.Apply(common.DefaultNamespace, common.Resource, originParent.UUID, originParent, false); err != nil {
				return nil, false, err
			}
		}

		parent.Spec.CpuRemain -= resource.Cpu
		parent.Spec.MemoryRemain -= resource.Memory

		parent.GenerateVersion()
		if _, _, err := a.IService.Apply(common.DefaultNamespace, common.Resource, parent.UUID, parent, false); err != nil {
			return nil, false, err
		}

		res.Spec.Cpu = resource.Cpu
		res.Spec.Memory = resource.Memory
		res.Spec.Pod = resource.Pod
		res.Spec.ParentApp = parent.UUID
		res.Spec.App = resource.App
		res.Metadata.Name = resource.Name

		res.GenerateVersion()
		if _, _, err := a.IService.Apply(common.DefaultNamespace, common.Resource, res.UUID, res, false); err != nil {
			return nil, false, err
		}

		updateResource = append(updateResource, res)

		// create history
		app := &appproject.AppProject{}
		if err := a.IService.GetByUUID(common.DefaultNamespace, common.Namespace, parent.Spec.App, app); err != nil {
			return nil, false, errors.New("the Namespace is not found")
		}

		// Get Creator

		history := &appproject.ConfigHistory{
			Spec: appproject.HistorySpec{
				App: res.Spec.App,
				History: map[string]interface{}{
					"creator":        "",
					"name":           res.Name,
					"namespace":      app.Spec.Desc,
					"pod":            res.Spec.Pod,
					"cpu_limit":      parent.Spec.Cpu,
					"memory_limit":   parent.Spec.Memory,
					"cpu_require":    res.Spec.Cpu,
					"memory_require": res.Spec.Memory,
				},
			},
		}

		if _, err := a.IService.Create(common.DefaultNamespace, common.History, history); err != nil {
			return nil, false, errors.New("the history create failed")
		}
	}

	// update config
	app := &appproject.AppProject{}
	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppProject, data.App, app); err != nil {
		return nil, false, errors.New("The app is not exist")
	}

	if app.Spec.AppType != appproject.App {
		return nil, false, errors.New("This is not an App type")
	}

	dbObj := &appproject.AppConfig{}
	a.IService.GetByFilter(common.DefaultNamespace, common.AppConfig, dbObj, map[string]interface{}{
		"spec.app": app.Metadata.UUID,
	})

	dbObj.Spec.Config = data.Config
	dbObj.Spec.App = app.Metadata.UUID

	dbObj.GenerateVersion()
	_, _, err = a.IService.Apply(common.DefaultNamespace, common.AppConfig, dbObj.UUID, dbObj, false)
	if err != nil {
		return nil, false, err
	}

	result := &apiResource.AppConfigResponse{
		Config:    dbObj,
		Resources: updateResource,
	}

	return result, true, nil
}

func (a *AppConfigService) DeleteResource(uuid string) error {
	return a.IService.Delete(common.DefaultNamespace, common.Resource, uuid)
}

func (a *AppConfigService) History(appid string, page, pageSize int64) ([]interface{}, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{
		"spec.app": appid,
	}

	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}
	return a.IService.ListByFilter(common.DefaultNamespace, common.History, filter, sort, offset, pageSize)
}
