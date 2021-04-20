package appservice

import (
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appproject"
	"github.com/yametech/devops/pkg/service"
)

type AppConfigService struct {
	service.IService
}

func NewAppConfigService(i service.IService) *AppConfigService {
	return &AppConfigService{i}
}

func (a *AppConfigService) GetByFilter(req *appproject.AppConfig) error {
	if err := a.IService.GetByFilter(common.DefaultNamespace, common.AppConfig, req, map[string]interface{}{
		"spec.app": req.Spec.App,
	}); err != nil {
		return err
	}

	return nil
}

func (a *AppConfigService) Update(req *appproject.AppConfig) (core.IObject, bool, error) {
	app := &appproject.AppProject{}
	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppProject, req.Spec.App, app); err != nil {
		return nil, false, errors.New("The app is not exist")
	}

	if app.Spec.AppType != appproject.App{
		return nil, false, errors.New("This is not an App type")
	}

	dbObj := &appproject.AppConfig{}
	a.IService.GetByFilter(common.DefaultNamespace, common.AppConfig, dbObj, map[string]interface{}{
		"spec.app": app.Metadata.UUID,
	})

	dbObj.Spec.Config = req.Spec.Config
	dbObj.Spec.App = app.Metadata.UUID

	dbObj.GenerateVersion()
	return a.IService.Apply(common.DefaultNamespace, common.AppConfig, dbObj.UUID, dbObj)
}
