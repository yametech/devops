package resource

import "github.com/yametech/devops/pkg/resource/appproject"

type AppProjectResponse struct {
	appproject.AppProject
	Children []*AppProjectResponse `json:"children"`
}
