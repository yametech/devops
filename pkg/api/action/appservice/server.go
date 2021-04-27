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
	*appService.NamespaceConfigService
}

func NewAppServiceServer(serviceName string, server *api.Server) *Server {
	cfaServer := &Server{
		Server:            server,
		AppProjectService: appService.NewAppProjectService(server.IService),
		AppConfigService: appService.NewAppConfigService(server.IService),
		NamespaceService: appService.NewNamespaceService(server.IService),
		NamespaceConfigService: appService.NewNamespaceConfigService(server.IService),
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
	}

	// Namespace
	{
		group.GET("/namespace", cfaServer.ListNamespaces)
		group.POST("/namespace", cfaServer.CreateNamespace)
		group.GET("/namespace/all", cfaServer.ListByLevel)
	}

	// NamespaceConfig
	{
		group.GET("/namespaceconfig/:uuid", cfaServer.GetNamespaceConfig)
		group.POST("/namespaceconfig", cfaServer.UpdateNamespaceConfig)
	}

	return cfaServer
}
