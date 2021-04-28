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
	*appService.ResourcePoolService
	*appService.ResourcePoolConfigService
}

func NewAppServiceServer(serviceName string, server *api.Server) *Server {
	cfaServer := &Server{
		Server:            server,
		AppProjectService: appService.NewAppProjectService(server.IService),
		AppConfigService: appService.NewAppConfigService(server.IService),
		ResourcePoolService: appService.NewResourcePoolService(server.IService),
		ResourcePoolConfigService: appService.NewResourcePoolConfigService(server.IService),
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
		group.GET("/app-config/:uuid", cfaServer.GetAppConfig)
		group.POST("/app-config", cfaServer.UpdateAppConfig)
		group.GET("/history/:uuid", cfaServer.ConfigHistory)
		group.DELETE("/app-config/resource/:uuid", cfaServer.DeleteResource)
	}

	// ResourcePool
	{
		group.GET("/resource-pool", cfaServer.ListResourcePool)
		group.POST("/resource-pool", cfaServer.CreateResourcePool)
		group.GET("/menu", cfaServer.ListByLevel)
	}

	// ResourcePoolConfig
	{
		group.GET("/resource-pool-config/:uuid", cfaServer.GetResourcePoolConfig)
		group.POST("/resource-pool-config", cfaServer.UpdateResourcePoolConfig)
	}

	return cfaServer
}


