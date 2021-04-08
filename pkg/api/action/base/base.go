package base

import (
	"fmt"
	"github.com/yametech/devops-zpk-server/pkg/api"
	serviceUser "github.com/yametech/devops-zpk-server/pkg/service/user"
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
		group.GET("/createUser", base.CreateUser)
	}

	//User Role
	{
		//group.GET("/user-role, base.ListRole)
	}

	return base
}
