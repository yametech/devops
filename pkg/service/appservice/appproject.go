package appservice

import (
	"encoding/json"
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appservice"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appservice"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type AppProjectService struct {
	service.IService
}

func NewAppProjectService(i service.IService) *AppProjectService {
	return &AppProjectService{i}
}

func (a *AppProjectService) List(search string) ([]*apiResource.Response, error) {
	if search != "" {
		return a.Search(search, 2)
	}

	sort := map[string]interface{}{
		"metadata.created_time": 1,
	}

	// Get the BusinessLine
	businessLine := &apiResource.Response{}
	if err := a.Children(businessLine, sort); err != nil {
		return nil, err
	}

	return businessLine.Children, nil
}

func (a *AppProjectService) Create(request *apiResource.Request) (core.IObject, error) {

	req := &appservice.AppProject{
		Metadata: core.Metadata{
			Name: request.Name,
		},
		Spec: appservice.AppSpec{
			AppType:   request.AppType,
			ParentApp: request.ParentApp,
			Desc:      request.Desc,
			Owner:     request.Owner,
		},
	}

	if req.Metadata.Name == "" {
		return nil, errors.New("The Name is requried")
	}

	filter := map[string]interface{}{
		"metadata.name": req.Name,
	}

	if err := a.IService.GetByFilter(common.DefaultNamespace, common.AppProject, req, filter); err == nil {
		return nil, errors.New("The Name is exist")
	}

	req.GenerateVersion()
	parent := &appservice.AppProject{}
	if req.Spec.ParentApp != "" {
		if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppProject, req.Spec.ParentApp, parent); err != nil {
			return nil, err
		}

		if parent.Spec.RootApp != "" {
			req.Spec.RootApp = parent.Spec.RootApp
		} else {
			req.Spec.RootApp = parent.Metadata.UUID
		}
	}
	return a.IService.Create(common.DefaultNamespace, common.AppProject, req)
}

func (a *AppProjectService) Update(uuid string, request *apiResource.Request) (core.IObject, bool, error) {
	req := &appservice.AppProject{
		Spec: appservice.AppSpec{
			Owner: request.Owner,
			Desc:  request.Desc,
		},
	}

	dbObj := &appservice.AppProject{}
	if err := a.IService.GetByUUID(common.DefaultNamespace, common.AppProject, uuid, dbObj); err != nil {
		return nil, false, err
	}
	if dbObj.UUID == "" {
		return nil, false, errors.New("The uuid is not exist")
	}

	dbObj.Spec.Desc = req.Spec.Desc
	dbObj.Spec.Owner = req.Spec.Owner

	return a.IService.Apply(common.DefaultNamespace, common.AppProject, uuid, dbObj, false)
}

func (a *AppProjectService) Delete(uuid string) (bool, error) {
	dbObj := &appservice.AppProject{}
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

func (a *AppProjectService) Children(req *apiResource.Response, sort map[string]interface{}) error {
	filter := map[string]interface{}{
		"spec.parent_app": req.UUID,
	}

	data, err := a.IService.ListByFilter(common.DefaultNamespace, common.AppProject, filter, sort, 0, 0)
	children := make([]*apiResource.Response, 0)
	if err = utils.UnstructuredObjectToInstanceObj(data, &children); err != nil {
		return err
	}

	if req.Spec.AppType == appservice.Service {
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

func (a *AppProjectService) Search(search string, level int64) ([]*apiResource.Response, error) {
	parentsMap := make(map[string]*apiResource.Response, 0)
	parents := make([]*apiResource.Response, 0)
	filter := make(map[string]interface{}, 0)

	sort := map[string]interface{}{
		"metadata.created_time": 1,
	}

	for ; level >= 0; level-- {

		filter["$or"] = []map[string]interface{}{
			{
				"metadata.name": bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
			},
			{
				"spec.desc":     bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
			},
		}

		filter["spec.app_type"] = level
		apps, err := a.IService.ListByFilter(common.DefaultNamespace, common.AppProject, filter, sort, 0, 0)
		if err != nil {
			continue
		}

		data := make([]*apiResource.Response, 0)
		if err = utils.UnstructuredObjectToInstanceObj(apps, &data); err != nil {
			return nil, err
		}

		// Get Root app
		for _, app := range data {
			if app.Spec.ParentApp == "" {
				if _, ok := parentsMap[app.Metadata.UUID]; !ok {
					parents = append(parents, app)
				}
			}

			if _, ok := parentsMap[app.Spec.RootApp]; app.Spec.RootApp != "" && !ok {
				root := &appservice.AppProject{}
				if err = a.IService.GetByUUID(common.DefaultNamespace, common.AppProject, app.Spec.RootApp, root); err != nil {
					continue
				}

				rootResponse := &apiResource.Response{}
				if err = utils.UnstructuredObjectToInstanceObj(root, &rootResponse); err != nil {
					return nil, err
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
			return nil, err
		}
	}
	return parents, nil
}


func (a *AppProjectService) SyncFromCMDB() ([]apiResource.CMDBData, error) {
	req := utils.NewRequest(http.Client{Timeout: 30 * time.Second}, "http", "cmdb-api-test.compass.ym", map[string]string{
		"Content-Type": "application/json",
	})
	resp, err := req.Get("/cmdb/api/v1/app-tree", nil)
	if err != nil {
		return nil, err
	}

	respData := make(map[string]interface{}, 0)
	if err = json.Unmarshal(resp, &respData); err != nil {
		return nil, err
	}

	result := make([]apiResource.CMDBData, 0)
	if data, ok := respData["data"]; ok {
		if err = utils.UnstructuredObjectToInstanceObj(data, &result); err != nil {
			return nil, err
		}
	}

	checkUpdate := apiResource.CMDBData{
		Children: result,
	}

	if err = a.UpdateFromCMDB(checkUpdate, ""); err != nil{
		return nil, err
	}



	return result, nil
}

func (a *AppProjectService) UpdateFromCMDB(data apiResource.CMDBData, parent string) error {
	if len(data.Children) <= 0{
		return nil
	}

	//parentObj := &appservice.AppProject{}
	//a.IService.GetByFilter(common.DefaultNamespace, common.AppProject, parentObj)
	//
	//for _, child := range data.Children{
	//
	//}

	return nil
}