package appservice

import "github.com/yametech/devops/pkg/service"

type AppProjectService struct {
	service.IService
}

func NewAppProjectService(i service.IService) *AppProjectService {
	return &AppProjectService{i}
}
