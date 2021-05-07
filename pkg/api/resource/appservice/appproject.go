package appservice

import (
	"github.com/yametech/devops/pkg/resource/appservice"
)

type Response struct {
	appservice.AppProject
	Children []*Response `json:"children"`
}

type Request struct {
	UUID      string             `json:"uuid"`
	Name      string             `json:"name"`
	AppType   appservice.AppType `json:"app_type"`
	ParentApp string             `json:"parent_app"`
	Desc      string             `json:"desc"`
	Owner     []string           `json:"owner"`
	Cpu       float64            `json:"cpu"`
	Memory    float64            `json:"memory"`
}
