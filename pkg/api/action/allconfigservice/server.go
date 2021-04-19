package allconfigservice

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	allConfigService "github.com/yametech/devops/pkg/service/allconfigservice"
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
		group.GET("/allConfig-project", allServer.ListAllConfigProject)
		group.POST("/allConfig-project", allServer.CreateAllConfigProject)
		group.PUT("/allConfig-project/:uuid", allServer.UpdateAllConfigProject)
		group.DELETE("/allConfig-project/:uuid", allServer.DeleteAllConfigProject)
	}
	return allServer
}
