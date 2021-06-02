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

// parse from cmdb data
type CMDBData struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	Leader   string     `json:"leader"`
	Desc     string     `json:"desc"`
	Owner    string     `json:"owner"`
	Children []CMDBData `json:"children"`
}
