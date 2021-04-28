package appservice

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appservice"
	"github.com/yametech/devops/pkg/resource/workorder"
	"github.com/yametech/devops/pkg/service"
	"log"
)

type AppConfigService struct {
	service.IService
}

func NewAppConfigService(i service.IService) *AppConfigService {
	return &AppConfigService{IService: i}
}

func (a *AppConfigService) GetByFilter(appid string) (core.IObject, error) {
	req := &appservice.AppConfig{}

	if err := a.IService.GetByFilter(common.DefaultNamespace, common.AppConfig, req, map[string]interface{}{
		"spec.app": appid,
	}); err != nil {
		return nil, err
	}
	return req, nil
}

func (a *AppConfigService) UpdateConfigResource(data *apiResource.ResourcePoolRequest) (core.IObject, bool, error) {

	appResource := appservice.AppResource{}

	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppResource, data.UUID, appResource); err != nil {
		return nil, false, err
	}

	// check the order status
	if appResource.Spec.OrderStatus == workorder.Checking {
		return nil, false, errors.New("the workorder is checking")
	}

	// Get Cpus and Memories from cmdb
	//cmdbCpus := 1000.0
	//cmdbMemory := 10240000

	// check Cpus and Memory

	// check the other resources


	// update the Resources Config


	// create history
	//history := appservice.AppConfigHistory{
	//	Spec: appservice.AppConfigHistorySpec{
	//		App: appResource.UUID,
	//	},
	//}

	// update config
	return nil, false, nil
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
	dbObj.GenerateVersion()
	return a.IService.Apply(common.DefaultNamespace, common.AppConfig, dbObj.UUID, dbObj, false)

}

func (a *AppConfigService) DeleteResource(uuid string) error {
	return a.IService.Delete(common.DefaultNamespace, common.AppResource, uuid)
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
