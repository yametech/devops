package appproject

import (
	"github.com/yametech/devops/pkg/resource/appproject"
)

type Response struct {
	appproject.AppProject
	Children []*Response `json:"children"`
}

type Request struct {
	Name      string             `json:"name"`
	AppType   appproject.AppType `json:"app_type"`
	ParentApp string             `json:"parent_app"`
	Desc      string             `json:"desc"`
	Owner     []string           `json:"owner"`
}
