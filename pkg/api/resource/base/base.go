package base

import (
	"github.com/yametech/devops/pkg/resource/base"
)

type ModuleRequest struct {
	UUID    string                 `json:"uuid"`
	Name    string                 `json:"name"`
	Parent  string                 `json:"parent"`
	Extends map[string]interface{} `json:"extends"`
}

type ModuleResponse struct {
	base.Module
	Children []*base.Module `json:"children"`
}
