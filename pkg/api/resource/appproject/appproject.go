package appproject

import (
	"github.com/yametech/devops/pkg/resource/appproject"
)

type AppProjectResponse struct {
	appproject.AppProject
	Children []*AppProjectResponse `json:"children"`
}

type AppProjectRequest struct {
	Name      string             `json:"name"`
	AppType   appproject.AppType `json:"app_type"`
	ParentApp string             `json:"parent_app"`
	Desc      string             `json:"desc"`
	Owner     []string           `json:"owner"`
}
