package appservice

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	appService "github.com/yametech/devops/pkg/service/appservice"
)

type Server struct {
	*api.Server
	*appService.AppProjectService
}

func NewAppServiceServer(serviceName string, server *api.Server) *Server {
	cfaServer := &Server{
		Server:            server,
		AppProjectService: appService.NewAppProjectService(server.IService),
	}
	group := cfaServer.Group(fmt.Sprintf("/%s", serviceName))

	//AppProject
	{
		group.GET("/app-project", cfaServer.ListAppProject)
		group.POST("/app-project", cfaServer.CreateAppProject)
		group.GET("/app-project/:uuid", cfaServer.RetrieveAppProject)
		group.PUT("/app-project/:uuid", cfaServer.UpdateAppProject)
		group.DELETE("/app-project/:uuid", cfaServer.DeleteAppProject)
	}

	return cfaServer
}
