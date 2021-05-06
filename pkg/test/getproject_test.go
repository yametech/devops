package test

import (
	"encoding/json"
	"fmt"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appservice"
	"github.com/yametech/devops/pkg/resource/workorder"
	"github.com/yametech/devops/pkg/store/mongo"
	"github.com/yametech/devops/pkg/utils"
	"io/ioutil"
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

func TestGetAppProject(t *testing.T) {
	b, err := ioutil.ReadFile("data.json") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	datas := make([]*BusinessLine, 0)
	json.Unmarshal(b, &datas)

	store, _, _ := mongo.NewMongo("mongodb://127.0.0.1:27017/admin")
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

func TestGenerateNumber(t *testing.T) {
	w := &workorder.WorkOrder{
		Spec: workorder.Spec{
			OrderType: 0,
		},
	}
	w.GenerateNumber()
}

func TestRequest(t *testing.T) {
	url := fmt.Sprintf("http://127.0.0.1:8081/workorder/status?relation=%s&order_type=%d",
		"57a093fb-d7fe-4875-b764-8da053994531", 1)
	body, _ := utils.Request("GET",
		url, nil, nil)

	fmt.Println(body)
	data := make(map[string]interface{})
	json.Unmarshal(body, &data)
	fmt.Println(data)
	fmt.Println(data["data"])
}
