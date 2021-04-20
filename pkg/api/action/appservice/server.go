package appservice

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	appService "github.com/yametech/devops/pkg/service/appservice"
)

type Server struct {
	*api.Server
	*appService.AppProjectService
	*appService.AppConfigService
}

func NewAppServiceServer(serviceName string, server *api.Server) *Server {
	cfaServer := &Server{
		Server:            server,
		AppProjectService: appService.NewAppProjectService(server.IService),
		AppConfigService: appService.NewAppConfigService(server.IService),
	}
	group := cfaServer.Group(fmt.Sprintf("/%s", serviceName))

	// AppProject
	{
		group.GET("/app-project", cfaServer.ListAppProject)
		group.POST("/app-project", cfaServer.CreateAppProject)
		group.PUT("/app-project/:uuid", cfaServer.UpdateAppProject)
		group.DELETE("/app-project/:uuid", cfaServer.DeleteAppProject)
	}

	// AppConfig
	{
		group.POST("/app-config", cfaServer.GetAppConfig)
		group.PUT("/app-config", cfaServer.UpdateAppConfig)
	}

	return cfaServer
}
