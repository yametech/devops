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
	*appService.NamespaceService
	*appService.ResourcePoolConfigService
}

func NewAppServiceServer(serviceName string, server *api.Server) *Server {
	cfaServer := &Server{
		Server:                    server,
		AppProjectService:         appService.NewAppProjectService(server.IService),
		AppConfigService:          appService.NewAppConfigService(server.IService),
		NamespaceService:          appService.NewResourcePoolService(server.IService),
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
		group.GET("/app-resource/:uuid", cfaServer.GetAppResource)
		group.POST("/app-resource", cfaServer.UpdateAppResource)
		group.DELETE("/app-resource/:uuid", cfaServer.DeleteResource)
	}

	// Namespace
	{
		group.GET("/namespace", cfaServer.ListNamespace)
		group.POST("/namespace", cfaServer.CreateNamespace)
		group.GET("/namespace/:uuid", cfaServer.GetNamespaceResourceRemain)
	}

	// ResourcePoolConfig
	{
		group.GET("/resource-pool-config/:uuid", cfaServer.GetResourcePoolConfig)
		group.POST("/resource-pool-config", cfaServer.UpdateResourcePoolConfig)
	}

	// Menu by level
	{
		group.GET("/menu", cfaServer.ListByLevel)
	}

	// History
	{
		group.GET("/history/:uuid", cfaServer.ConfigHistory)
	}

	return cfaServer
}
