package appproject

import (
	"github.com/yametech/devops/pkg/resource/appservice"
)

type Response struct {
	appservice.AppProject
	Children []*Response `json:"children"`
}

type Request struct {
	Name      string             `json:"name"`
	AppType   appservice.AppType `json:"app_type"`
	ParentApp string             `json:"parent_app"`
	Desc      string             `json:"desc"`
	Owner     []string           `json:"owner"`
}
