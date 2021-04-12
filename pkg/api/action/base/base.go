package base

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	baseService "github.com/yametech/devops/pkg/service/base"
)

type baseServer struct {
	*api.Server
	*baseService.UserService
	*baseService.UserProjectService
}

func NewBaseServer(serviceName string, server *api.Server) *baseServer {
	base := &baseServer{
		Server:             server,
		UserService:        baseService.NewUser(server.IService),
		UserProjectService: baseService.NewUserProjectService(server.IService),
	}
	group := base.Group(fmt.Sprintf("/%s", serviceName))

	//UserProjectService
	{
		group.GET("/users", base.ListUser)
		group.GET("/user/:uuid", base.GetUser)
		group.POST("/user", base.CreateUser)
		group.PUT("/user/:uuid", base.UpdateUser)
		group.DELETE("/user/:uuid", base.DeleteUser)
	}

	// UserProject
	{
		group.POST("/project", base.CreateProject)
		group.GET("/project", base.ListProject)

	}

	return base
}
