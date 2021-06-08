package test

import (
	"encoding/json"
	"fmt"
	apiResource "github.com/yametech/devops/pkg/api/resource/appservice"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appservice"
	ser "github.com/yametech/devops/pkg/service/appservice"
	"github.com/yametech/devops/pkg/resource/workorder"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/store/mongo"
	"github.com/yametech/devops/pkg/utils"
	"io/ioutil"
	"reflect"
	"testing"
)

type BusinessLine struct {
	Serid    string     `json:"ser_id"`
	Name     string     `json:"name"`
	Children []*Service `json:"children"`
}

type Service struct {
	Busid    string `json:"bus_id"`
	Name     string `json:"name"`
	Children []*App `json:"children"`
}

type App struct {
	Appid string   `json:"app_id"`
	Name  string   `json:"name"`
	Desc  string   `json:"desc"`
	Owner []string `json:"owner"`
}

type NamespaceBusinessLine struct {
	Id           string           `json:"id"`
	BusinessName string           `json:"business_name"`
	Leaders      string           `json:"leader"`
	Children     []*NamespaceLine `json:"children"`
}

type NamespaceLine struct {
	Id        string                 `json:"id"`
	Namespace string                 `json:"namespace"`
	Env       string                 `json:"env"`
	Config    map[string]interface{} `json:"config"`
	Threshold int                    `json:"threshold"`
}

func TestGetAppProject(t *testing.T) {
	b, err := ioutil.ReadFile("appproejct.json") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	datas := make([]*BusinessLine, 0)
	json.Unmarshal(b, &datas)

	store, _, _ := mongo.NewMongo("mongodb://10.200.10.46:27017/devops")
	for _, data := range datas {
		buinessLine := &appservice.AppProject{
			Metadata: core.Metadata{
				Name: data.Name,
			},
			Spec: appservice.AppSpec{
				ParentApp: "",
				RootApp:   "",
				AppType:   appservice.BusinessLine,
				Desc:      "",
				Owner:     nil,
			},
		}
		store.Create(common.DefaultNamespace, common.AppProject, buinessLine)
		for _, services := range data.Children {
			service := &appservice.AppProject{
				Metadata: core.Metadata{
					Name: services.Name,
				},
				Spec: appservice.AppSpec{
					ParentApp: buinessLine.UUID,
					RootApp:   buinessLine.UUID,
					AppType:   appservice.Service,
					Desc:      "",
					Owner:     nil,
				},
			}
			store.Create(common.DefaultNamespace, common.AppProject, service)
			for _, apps := range services.Children {
				app := &appservice.AppProject{
					Metadata: core.Metadata{
						Name: apps.Desc,
					},
					Spec: appservice.AppSpec{
						ParentApp: service.UUID,
						RootApp:   service.Spec.RootApp,
						AppType:   appservice.App,
						Desc:      apps.Name,
						Owner:     apps.Owner,
					},
				}
				store.Create(common.DefaultNamespace, common.AppProject, app)
			}
		}
	}

	fmt.Println("success")
}

func TestGetNamespace(t *testing.T) {
	b, err := ioutil.ReadFile("namespace.json") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	datas := make([]*NamespaceBusinessLine, 0)
	json.Unmarshal(b, &datas)

	store, _, _ := mongo.NewMongo("mongodb://10.200.10.46:27017/devops")
	baseService := service.NewBaseService(store)
	namespaceService := ser.NewNamespaceService(baseService)

	for _, data := range datas {
		buinessLine := &appservice.AppProject{}
		filter := map[string]interface{}{
			"metadata.name": data.BusinessName,
		}
		store.GetByFilter(common.DefaultNamespace, common.AppProject, buinessLine, filter)
		buinessLine.Spec.Owner = append(buinessLine.Spec.Owner, data.Leaders)
		store.Apply(common.DefaultNamespace, common.AppProject, buinessLine.UUID, buinessLine, true)

		for _, child := range data.Children{
			namespace := &apiResource.Request{
				Name: child.Env,
				Desc: child.Namespace,
				ParentApp: buinessLine.UUID,
			}
			namespaceService.Update(namespace)
		}
	}

	fmt.Println("success")
}

func TestGenerateNumber(t *testing.T) {
	w := &workorder.WorkOrder{
		Spec: workorder.Spec{
			OrderType: 0,
		},
	}
	w.GenerateNumber()
}

func TestRequest(t *testing.T) {
	store, _, _ := mongo.NewMongo("mongodb://127.0.0.1:27017/devops")
	baseService := service.NewBaseService(store)
	service := ser.NewAppProjectService(baseService)
	service.SyncFromCMDB()
}

func TestCleanData(t *testing.T) {
	store, _, _ := mongo.NewMongo("mongodb://10.200.10.46:27017/devops")
	//store, _, _ := mongo.NewMongo("mongodb://127.0.0.1:27017/devops")
	appConfigs := make([]*appservice.AppConfig, 0)
	data, _ := store.ListByFilter(common.DefaultNamespace, common.AppConfig, nil, nil, 0, 0)
	utils.UnstructuredObjectToInstanceObj(data, &appConfigs)
	appsMap := make(map[string]*appservice.AppProject, 0)
	for _, appConfig := range appConfigs {
		app := &appservice.AppProject{}
		store.GetByUUID(common.DefaultNamespace, common.AppProject, appConfig.Spec.App, app)
		if app.UUID != ""{
			appsMap[app.UUID] = app
		}
	}

	apps := make([]*appservice.AppProject, 0)
	filter := map[string]interface{}{
		"spec.app_type": 2,
	}
	datas, _ := store.ListByFilter(common.DefaultNamespace, common.AppProject, filter, nil, 0, 0)
	utils.UnstructuredObjectToInstanceObj(datas, &apps)

	for _, app := range apps{
		if _, ok := appsMap[app.UUID]; !ok{
			store.Delete(common.DefaultNamespace, common.AppProject, app.UUID)
		}
	}
}

func TestEqual(t *testing.T) {
	m1 := map[string]interface{}{"1": []int{1, 2, 3}, "2": 3, "3": "a", "4": map[int]interface{}{1: 1, 2: 2}}
	m2 := map[string]interface{}{"1": []int{1, 2, 3}, "2": 3, "3": "a", "4": map[int]interface{}{1: 1, 2: 2}}
	if reflect.DeepEqual(m1, m2) {
		fmt.Println("相等")
	}
}

func TestName(t *testing.T) {
	ss := make([]int, 0)
	ss = append(ss, 1)
	ss = nil
	ss = append(ss, 2)
}


