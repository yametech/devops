package appservice

import (
	"fmt"
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appproject"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
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
		Config: req,
		Resources: resource,
	}

	return &response, nil
}

func (a *AppConfigService) Update(data *apiResource.AppConfigRequest) (core.IObject, bool, error) {

	app := &appproject.AppProject{}
	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppProject, data.App, app); err != nil {
		return nil, false, errors.New("The app is not exist")
	}

	if app.Spec.AppType != appproject.App {
		return nil, false, errors.New("This is not an App type")
	}

	// Merging request resources from the same namespace


	// check the resource
	for _, resource := range data.Resources {
		parent := &appproject.Resource{}
		if err := a.IService.GetByUUID(common.DefaultNamespace, common.Resource, resource.ParentApp, parent); err != nil {
			return nil, false, errors.New("The resource is not exist")
		}
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
		limitCpu := int(((parent.Spec.Cpus - (parent.Spec.CpuRemain - resource.Cpu)) / parent.Spec.Cpus) * 100)
		if limitCpu > parent.Spec.Threshold{
			err := fmt.Sprintf("CPU total resources exceed %d", parent.Spec.Threshold)
			return nil, false, errors.New(err)
		}
		limitMemory := int(((parent.Spec.Memories - (parent.Spec.MemoryRemain - resource.Memory)) / parent.Spec.Memory) * 100)
		if limitMemory > parent.Spec.Threshold{
			err := fmt.Sprintf("Memory total resources exceed %d", parent.Spec.Threshold)
			return nil, false, errors.New(err)
		}
	}

	for _, resource := range data.Resources {
		parent := &appproject.Resource{}
		if err := a.IService.GetByUUID(common.DefaultNamespace, common.Resource, resource.ParentApp, parent); err != nil {
			return nil, false, errors.New("The resource is not exist")
		}

		res := &appproject.Resource{}
		if err := a.IService.GetByUUID(common.DefaultNamespace, common.Resource, resource.App, res); err != nil {
			return nil, false, errors.New("The resource is not exist")
		}

		originParent := &appproject.Resource{}
		if err := a.IService.GetByUUID(common.DefaultNamespace, common.Resource, res.Spec.ParentApp, originParent); err != nil {
			return nil, false, errors.New("The resource is not exist")
		}

		originParent.Spec.CpuRemain += res.Spec.Cpu
		originParent.Spec.MemoryRemain += res.Spec.Memory

		originParent.GenerateVersion()
		if _, _, err := a.IService.Apply(common.DefaultNamespace, common.Resource, originParent.UUID, originParent, false); err != nil {
			return nil, false, err
		}

		parent.Spec.CpuRemain -= resource.Cpu
		parent.Spec.MemoryRemain -= resource.Memory

		parent.GenerateVersion()
		if _, _, err := a.IService.Apply(common.DefaultNamespace, common.Resource, parent.UUID, parent, false); err != nil {
			return nil, false, err
		}

		res.Spec.Cpu = resource.Cpu
		res.Spec.Memory = resource.Memory
		res.Spec.ParentApp = parent.UUID

		res.GenerateVersion()
		if _, _, err := a.IService.Apply(common.DefaultNamespace, common.Resource, res.UUID, res, false); err != nil {
			return nil, false, err
		}
	}

	dbObj := &appproject.AppConfig{}
	a.IService.GetByFilter(common.DefaultNamespace, common.AppConfig, dbObj, map[string]interface{}{
		"spec.app": app.Metadata.UUID,
	})

	dbObj.Spec.Config = data.Config
	dbObj.Spec.App = app.Metadata.UUID

	dbObj.GenerateVersion()
	return a.IService.Apply(common.DefaultNamespace, common.AppConfig, dbObj.UUID, dbObj, false)
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