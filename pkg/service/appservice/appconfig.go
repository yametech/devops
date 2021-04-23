package appservice

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appproject"
	"github.com/yametech/devops/pkg/service"
)

type AppConfigService struct {
	service.IService
	UpdateMap map[appproject.AppType]func(a service.IService, data *apiResource.AppConfigRequest) (*appproject.AppProject, error)
}

func NewAppConfigService(i service.IService) *AppConfigService {
	updateMap := map[appproject.AppType]func(a service.IService, data *apiResource.AppConfigRequest) (*appproject.AppProject, error){
		appproject.App: ExecAppType,
		appproject.Namespace: ExecNameSpaceType,
	}
	return &AppConfigService{IService: i, UpdateMap: updateMap}
}


func (a *AppConfigService) GetByFilter(appid string) (core.IObject, error) {
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

	return req, nil
}

func (a *AppConfigService) Update(data *apiResource.AppConfigRequest) (core.IObject, bool, error) {

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
	return a.IService.Apply(common.DefaultNamespace, common.AppConfig, dbObj.UUID, dbObj, false)
}


func ExecAppType(a service.IService, data *apiResource.AppConfigRequest) (*appproject.AppProject, error) {
	app := &appproject.AppProject{}
	if err := a.GetByUUID(common.DefaultNamespace, common.AppProject, data.App, app); err != nil {
		return nil, errors.New("The app is not exist")
	}

	if app.Spec.AppType != appproject.App {
		return nil, errors.New("This is not an App type")
	}

	return app, nil
}

func ExecNameSpaceType(a service.IService, data *apiResource.AppConfigRequest) (*appproject.AppProject, error) {
	app := &appproject.AppProject{}
	if err := a.GetByUUID(common.DefaultNamespace, common.Namespace, data.App, app); err != nil {
		return nil, errors.New("The namespace is not exist")
	}

	if app.Spec.AppType != appproject.Namespace {
		return nil, errors.New("This is not an Namespace type")
	}

	return app, nil
}