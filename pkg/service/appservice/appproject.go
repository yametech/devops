package appservice

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"time"


	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appservice"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appservice"
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

func (a *AppProjectService) List(search string) ([]*apiResource.Response, error) {
	if search != "" {
		return a.Search(search, 2)
	}

	sort := map[string]interface{}{
		"metadata.create": 1,
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
				"spec.desc": bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
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

func (a *AppProjectService) SyncFromCMDB() error {
	req := utils.NewRequest(http.Client{Timeout: 30 * time.Second}, "http", "cmdb-api.compass.ym", map[string]string{
		"Content-Type": "application/json",
	})
	resp, err := req.Get("/cmdb/api/v1/app-tree", nil)
	if err != nil {
		return err
	}

	respData := make(map[string]interface{}, 0)
	if err = json.Unmarshal(resp, &respData); err != nil {
		return err
	}

	if code, ok := respData["code"]; ok {
		if code.(float64) != 200 {
			return fmt.Errorf("request cmdb response code: %v", code)
		}
	}

	result := make([]apiResource.CMDBData, 0)
	if data, ok := respData["data"]; ok {
		if err = utils.UnstructuredObjectToInstanceObj(data, &result); err != nil {
			return err
		}
	}

	if err = a.DeleteFromCMDB(result, map[string]interface{}{"spec.app_type": 0}, 0); err != nil {
		return err
	}

	for _, businessLine := range result {
		if err = a.UpdateBusinessFromCMDB(businessLine, "", ""); err != nil {
			return err
		}
	}

	return nil
}

func (a *AppProjectService) DeleteFromCMDB(cmdb []apiResource.CMDBData, filter map[string]interface{}, level int) error {

	dbProject := make([]*appservice.AppProject, 0)
	data, err := a.IService.ListByFilter(common.DefaultNamespace, common.AppProject, filter, nil, 0, 0)
	if err != nil {
		return err
	}

	if err = utils.UnstructuredObjectToInstanceObj(data, &dbProject); err != nil {
		return err
	}

	cmdbMap := make(map[string]struct{}, 0)

	for _, obj := range cmdb {
		if _, ok := cmdbMap[obj.Name]; !ok {
			cmdbMap[obj.Name] = struct{}{}
		}
	}

	for _, dbObj := range dbProject {
		if level == 2 {
			if _, ok := cmdbMap[dbObj.Spec.Desc]; !ok {
				if err = a.DeleteEveryLevel(dbObj); err != nil {
					return err
				}
			}
		} else {
			if _, ok := cmdbMap[dbObj.Name]; !ok {
				if err = a.DeleteEveryLevel(dbObj); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (a *AppProjectService) DeleteEveryLevel(req *appservice.AppProject) error {

	if err := a.IService.Delete(common.DefaultNamespace, common.AppProject, req.UUID); err != nil {
		return err
	}
	
  
  .Printf("[Controller] appproject delete: [Name: %v, Desc: %v]\n", req.Name, req.Spec.Desc)
	children := make([]*appservice.AppProject, 0)
	filter := map[string]interface{}{
		"spec.parent_app": req.UUID,
	}
	data, err := a.IService.ListByFilter(common.DefaultNamespace, common.AppProject, filter, nil, 0, 0)
	if err != nil {
		return nil
	}

	if err = utils.UnstructuredObjectToInstanceObj(data, &children); err != nil {
		return err
	}

	for _, child := range children {
		if err = a.DeleteEveryLevel(child); err != nil {
			return err
		}
	}

	return nil
}

func (a *AppProjectService) UpdateBusinessFromCMDB(data apiResource.CMDBData, parent string, root string) error {

	if data.Name == "" {
		return nil
	}

	dbObj := &appservice.AppProject{}
	filter := map[string]interface{}{
		"metadata.name":   data.Name,
		"spec.parent_app": parent,
		"spec.root_app":   root,
	}

	if err := a.IService.GetByFilter(common.DefaultNamespace, common.AppProject, dbObj, filter); err != nil {
		dbObj.Metadata.Name = data.Name
		dbObj.Spec.AppType = 0
		dbObj.Spec.Owner = []string{data.Leader}

		if _, err = a.IService.Create(common.DefaultNamespace, common.AppProject, dbObj); err != nil {
			return err
		}

		log.Printf("[Controller] appproject add new BusinessLine: %v\n", data.Name)
	} else {

		if data.Leader != "" && len(dbObj.Spec.Owner) > 0 {
			if dbObj.Spec.Owner[0] != data.Leader {
				dbObj.Spec.Owner = []string{data.Leader}
				if _, _, err = a.IService.Apply(common.DefaultNamespace, common.AppProject, dbObj.UUID, dbObj, false); err != nil {
					return err
				}
				log.Printf("[Controller] appproject update BusinessLine: %v---Leader: %v\n", dbObj.Name, dbObj.Spec.Owner)

			}
			log.Printf("[Controller] appproject update BusinessLine: %v---Leader: %v\n", dbObj.Name, dbObj.Spec.Owner)
		}
	}

	if err := a.DeleteFromCMDB(data.Children, map[string]interface{}{"spec.parent_app": dbObj.UUID, "spec.app_type": 1}, 1); err != nil {
		return err
	}

	for _, child := range data.Children {
		if err := a.UpdateServiceFromCMDB(child, dbObj.UUID, dbObj.UUID); err != nil {
			return err
		}
	}

	return nil
}

func (a *AppProjectService) UpdateServiceFromCMDB(data apiResource.CMDBData, parent string, root string) error {
	if data.Name == "" {
		return nil
	}

	dbObj := &appservice.AppProject{}
	filter := map[string]interface{}{
		"metadata.name":   data.Name,
		"spec.parent_app": parent,
		"spec.root_app":   root,
	}

	if err := a.IService.GetByFilter(common.DefaultNamespace, common.AppProject, dbObj, filter); err != nil {
		dbObj.Metadata.Name = data.Name
		dbObj.Spec.ParentApp = parent
		dbObj.Spec.RootApp = root
		dbObj.Spec.AppType = 1

		if _, err = a.IService.Create(common.DefaultNamespace, common.AppProject, dbObj); err != nil {
			return err
		}

		log.Printf("[Controller] appproject add new Service: %v\n", data.Name)
	}

	if err := a.DeleteFromCMDB(data.Children, map[string]interface{}{"spec.parent_app": dbObj.UUID, "spec.app_type": 2}, 2); err != nil {
		return err
	}

	for _, child := range data.Children {
		if err := a.UpdateAppFromCMDB(child, dbObj.UUID, root); err != nil {
			return err
		}
	}

	return nil
}

func (a *AppProjectService) UpdateAppFromCMDB(data apiResource.CMDBData, parent string, root string) error {
	if data.Desc == "" {
		return nil
	}

	dbObj := &appservice.AppProject{}
	filter := map[string]interface{}{
		"metadata.name":   data.Desc,
		"spec.parent_app": parent,
		"spec.root_app":   root,
	}

	if err := a.IService.GetByFilter(common.DefaultNamespace, common.AppProject, dbObj, filter); err != nil {
		dbObj.Metadata.Name = data.Desc
		dbObj.Spec.Desc = data.Name
		dbObj.Spec.ParentApp = parent
		dbObj.Spec.RootApp = root
		dbObj.Spec.AppType = 2
		dbObj.Spec.Owner = []string{data.Owner}

		if _, err = a.IService.Create(common.DefaultNamespace, common.AppProject, dbObj); err != nil {
			return err
		}
		log.Printf("[Controller] appproject add new App: %v\n", dbObj.Spec.Desc)
	} else {
		if data.Owner != "" && len(dbObj.Spec.Owner) > 0 {
			if dbObj.Spec.Owner[0] != data.Owner {
				dbObj.Spec.Owner = []string{data.Owner}
				if _, _, err = a.IService.Apply(common.DefaultNamespace, common.AppProject, dbObj.UUID, dbObj, false); err != nil {
					return err
				}
				log.Printf("[Controller] appproject update App: %v----Owner: %v\n", dbObj.Spec.Desc, dbObj.Spec.Owner)
			}
		}
	}

	return nil
}
