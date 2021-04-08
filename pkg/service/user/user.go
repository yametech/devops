package user

import "github.com/yametech/devops-zpk-server/pkg/service"

type User struct {
	service.IService
}

func NewUser(i service.IService) *User {
	return &User{i}
}



