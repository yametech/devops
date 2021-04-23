package globalconfigservice

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	allConfigService "github.com/yametech/devops/pkg/service/globalservice"
)

type Server struct {
	*api.Server
	*allConfigService.GlobalConfigService
}

func NewGlobalServiceServer(serviceName string, server *api.Server) *Server {
	allServer := &Server{
		server,
		allConfigService.NewAllConfigService(server.IService),
	}
	group := allServer.Group(fmt.Sprintf("/%s", serviceName))

	//allConfigProject
	{
		group.GET("/globalconfig-project", allServer.ListGlobalConfig)
		//group.POST("/globalconfig-project", allServer.CreateGlobalConfig)
		group.POST("/globalconfig-project", allServer.UpdateGlobalConfig)
	}
	return allServer
}
