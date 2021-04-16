package Allconfigservice

import "github.com/yametech/devops/pkg/service"

type Allconfigservice struct {
	service.IService
}

func NewAllconfigservice(i service.IService) *Allconfigservice {
	return &Allconfigservice{i}
}

func List() {

}
