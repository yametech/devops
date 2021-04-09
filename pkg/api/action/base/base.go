package base

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	serviceUser "github.com/yametech/devops/pkg/service/user"
)

type baseServer struct {
	*api.Server
	*serviceUser.User
}

func NewBaseServer(serviceName string, server *api.Server) *baseServer {

	base := &baseServer{
		Server: server,
		User:   serviceUser.NewUser(server.IService),
	}
	group := base.Group(fmt.Sprintf("/%s", serviceName))

	//User
	{
		group.GET("/users", base.ListUser)
		group.GET("/user", base.GetUser)
		group.POST("/user", base.CreateUser)
		group.PUT("/user/:uuid", base.UpdateUser)
		group.DELETE("user/:uuid", base.DeleteUser)
	}

	//User Role
	{
		//group.GET("/user-role, base.ListRole)
	}

	return base
}
