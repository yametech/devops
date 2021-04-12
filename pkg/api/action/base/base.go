package base

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	serviceUser "github.com/yametech/devops/pkg/service/base"
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
		group.GET("/base", base.GetUser)
		group.POST("/base", base.CreateUser)
		group.PUT("/base/:uuid", base.UpdateUser)
		group.DELETE("base/:uuid", base.DeleteUser)
	}

	// Artifact
	{
		group.GET("/artifact", base.ListArtifact)
	}

	return base
}
