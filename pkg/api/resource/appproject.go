package resource

import "github.com/yametech/devops/pkg/resource"

type AppProjectResponse struct {
	resource.AppProject
	Children []*AppProjectResponse `json:"children"`
}
