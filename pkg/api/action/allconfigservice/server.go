package allconfigservice

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	allConfigService "github.com/yametech/devops/pkg/service/globalservice"
)

type Server struct {
	*api.Server
	*allConfigService.AllConfigService
}

func NewAllServiceServer(serviceName string, server *api.Server) *Server {
	allServer := &Server{
		server,
		allConfigService.NewAllConfigService(server.IService),
	}
	group := allServer.Group(fmt.Sprintf("/%s", serviceName))

	//allConfigProject
	{
		group.GET("/allConfig-project", allServer.ListGlobalConfig)
		group.POST("/allConfig-project", allServer.CreateGlobalConfig)
		group.PUT("/allConfig-project/:uuid", allServer.UpdateGlobalConfig)
	}
	return allServer
}
