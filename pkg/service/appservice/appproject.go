package appservice

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appproject"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppProjectService struct {
	service.IService
}

func NewAppProjectService(i service.IService) *AppProjectService {
	return &AppProjectService{i}
}

func (a *AppProjectService) List(search string) ([]*apiResource.AppProjectResponse, int64, error) {
	if search != "" {
		return a.Search(search, 2)
	}

	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}

	// Get the BusinessLine
	businessLine := &apiResource.AppProjectResponse{}
	if err := a.Children(businessLine, sort); err != nil {
		return nil, 0, err
	}

	return businessLine.Children, int64(len(businessLine.Children)), nil
}

func (a *AppProjectService) Create(req *appproject.AppProject) error {
	if req.Metadata.Name == "" {
		return errors.New("The Name is requried")
	}

	filter := map[string]interface{}{
		"metadata.name": req.Name,
	}

	if err := a.IService.GetByFilter(common.DefaultNamespace, common.AppProject, req, filter); err == nil {
		return errors.New("The Name is exist")
	}

	req.GenerateVersion()
	parent := &appproject.AppProject{}
	if req.Spec.ParentApp != "" {
		if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppProject, req.Spec.ParentApp, parent); err != nil {
			return err
		}

		if parent.Spec.RootApp != "" {
			req.Spec.RootApp = parent.Spec.RootApp
		} else {
			req.Spec.RootApp = parent.Metadata.UUID
		}
	}
	_, err := a.IService.Create(common.DefaultNamespace, common.AppProject, req)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppProjectService) Update(uuid string, req *appproject.AppProject) (core.IObject, bool, error) {
	dbObj := &appproject.AppProject{}
	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppProject, uuid, dbObj); err != nil {
		return nil, false, err
	}
	if dbObj.UUID == "" {
		return nil, false, errors.New("The uuid is not exist")
	}

	dbObj.Spec.Desc = req.Spec.Desc
	dbObj.Spec.Owner = req.Spec.Owner

	dbObj.GenerateVersion()
	return a.IService.Apply(common.DefaultNamespace, common.AppProject, uuid, dbObj, false)
}

func (a *AppProjectService) Delete(uuid string) (bool, error) {
	dbObj := &appproject.AppProject{}
	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppProject, uuid, dbObj); err != nil {
		return false, err
	}
	filter := map[string]interface{}{
		"spec.parent_app": dbObj.Metadata.UUID,
	}
	children, err := a.IService.ListByFilter(common.DefaultNamespace, common.AppProject, filter, nil, 0, 0)
	if err != nil {
		return false, err
	}

	if len(children) > 0 {
		return false, errors.New("This label has children labels, Please Delete them first")
	}

	err = a.IService.Delete(common.DefaultNamespace, common.AppProject, uuid)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *AppProjectService) Children(req *apiResource.AppProjectResponse, sort map[string]interface{}) error {
	filter := map[string]interface{}{
		"spec.parent_app": req.UUID,
	}

	data, err := a.IService.ListByFilter(common.DefaultNamespace, common.AppProject, filter, sort, 0, 0)
	children := make([]*apiResource.AppProjectResponse, 0)
	if err = utils.UnstructuredObjectToInstanceObj(data, &children); err != nil {
		return err
	}

	if req.Spec.AppType == appproject.Service {
		req.Children = children
		return nil
	}

	for _, child := range children {
		_child := child
		if err = a.Children(_child, sort); err != nil {
			return err
		}
	}

	req.Children = children
	return nil
}

func (a *AppProjectService) Search(search string, level int64) ([]*apiResource.AppProjectResponse, int64, error) {
	parentsMap := make(map[string]*apiResource.AppProjectResponse, 0)
	parents := make([]*apiResource.AppProjectResponse, 0)
	filter := make(map[string]interface{}, 0)
	if search != "" {
		filter["metadata.name"] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}}
	}

	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}

	for ; level >= 0; level-- {

		filter["spec.app_type"] = level
		apps, err := a.IService.ListByFilter(common.DefaultNamespace, common.AppProject, filter, sort, 0, 0)
		if err != nil {
			continue
		}

		data := make([]*apiResource.AppProjectResponse, 0)
		if err = utils.UnstructuredObjectToInstanceObj(apps, &data); err != nil {
			return nil, 0, err
		}

		// Get Root app
		for _, app := range data {
			if app.Spec.ParentApp == "" {
				if _, ok := parentsMap[app.Metadata.UUID]; !ok {
					parents = append(parents, app)
				}
			}

			if _, ok := parentsMap[app.Spec.RootApp]; app.Spec.RootApp != "" && !ok {
				root := &appproject.AppProject{}
				if err = a.IService.GetByUUID(common.DefaultNamespace, common.AppProject, app.Spec.RootApp, root); err != nil {
					continue
				}

				rootResponse := &apiResource.AppProjectResponse{}
				if err = utils.UnstructuredObjectToInstanceObj(root, &rootResponse); err != nil {
					return nil, 0, err
				}
				parentsMap[app.Spec.RootApp] = rootResponse
				parents = append(parents, rootResponse)
			}
		}
	}

	// Get the children of BusinessLine
	for _, child := range parents {
		_child := child
		if err := a.Children(_child, sort); err != nil {
			return nil, 0, err
		}
	}
	return parents, int64(len(parents)), nil
}
