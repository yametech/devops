package test

import (
	"encoding/json"
	"fmt"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appproject"
	"github.com/yametech/devops/pkg/store/mongo"
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
		buinessLine := &appproject.AppProject{
			Metadata: core.Metadata{
				Name: data.Name,
			},
			Spec: appproject.AppSpec{
				ParentApp: "",
				RootApp:   "",
				AppType:   appproject.BusinessLine,
				Desc:      "",
				Owner:     nil,
			},
		}
		store.Create(common.DefaultNamespace, common.AppProject, buinessLine)
		for _, services := range data.Children {
			service := &appproject.AppProject{
				Metadata: core.Metadata{
					Name: services.Name,
				},
				Spec: appproject.AppSpec{
					ParentApp: buinessLine.UUID,
					RootApp:   buinessLine.UUID,
					AppType:   appproject.Service,
					Desc:      "",
					Owner:     nil,
				},
			}
			store.Create(common.DefaultNamespace, common.AppProject, service)
			for _, apps := range services.Children {
				app := &appproject.AppProject{
					Metadata: core.Metadata{
						Name: apps.Name,
					},
					Spec: appproject.AppSpec{
						ParentApp: service.UUID,
						RootApp:   service.Spec.RootApp,
						AppType:   appproject.App,
						Desc:      apps.Desc,
						Owner:     apps.Owner,
					},
				}
				store.Create(common.DefaultNamespace, common.AppProject, app)
			}
		}
	}

	fmt.Println("success")
}

func TestGetData(t *testing.T) {
}
