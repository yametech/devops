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
const (
	BusinessUuid = "business"
	ServiceUuid  = "business_domain"
	AppUuid      = "business_service"
)


type Data struct {
	TotalCount int                      `json:"totalCount"`
	List       []map[string]interface{} `json:"list"`
	Code       int                      `json:"code"`
	Msg        string                   `json:"msg"`
}

type Base struct {
	Id   string `json:"id"`
	UUID string `json:"uuid"`
}

type Business struct {
	Base
	BusinessDescribe string `json:"business_describe"`
	BusinessMaster   string `json:"business_master"`
	BusinessName     string `json:"business_name"`
}

type Service struct {
	Base
	DomainId      string `json:"domain_id"`
	DomainName    string `json:"domain_name"`
	DomainRemarks string `json:"domain_remarks"`
}

type App struct {
	Base
	ServiceDescribe string `json:"service_describe"`
	ServiceGrade    string `json:"service_grade"`
	ServiceId       string `json:"service_id"`
}