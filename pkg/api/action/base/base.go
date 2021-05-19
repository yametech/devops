package base

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	baseService "github.com/yametech/devops/pkg/service/base"
)

type baseServer struct {
	*api.Server
	*baseService.GlobalModuleService
	*baseService.CollectionModuleService
	*baseService.AllModuleService
	*baseService.ModuleEntry
	*baseService.RecentVisit
	*baseService.ShowAllGroupModule
}

func NewBaseServer(serviceName string, server *api.Server) *baseServer {
	base := &baseServer{
		Server:                  server,
		GlobalModuleService:     baseService.NewGlobalModuleService(server.IService),
		CollectionModuleService: baseService.NewCollectionModuleService(server.IService),
		AllModuleService:        baseService.NewAllModuleService(server.IService),
		ModuleEntry:             baseService.NewModuleEntry(server.IService),
		RecentVisit:             baseService.NewRecentVisit(server.IService),
		ShowAllGroupModule:      baseService.NewShowAllGroupModule(server.IService),
	}
	group := base.Group(fmt.Sprintf("/%s", serviceName))
	// globalmodule
	{
		group.GET("/globalmodule", base.ListGlobalModule)
		group.POST("/globalmodule", base.CreateGlobalModule)
		group.DELETE("/globalmodule/:uuid", base.DeleteGlobalModule)
	}

	// collectionmodule
	{
		group.GET("/collectionmodule", base.ListCollectionModule)
		group.POST("/collectionmodule", base.AddCollectionModule)
		group.DELETE("/collectionmodule", base.DeleteCollectionModule)
	}

	// allmodule
	{
		group.GET("/allmodule", base.ListAll)
		group.POST("/allmodule/group", base.CreateGroup)
		group.POST("/allmodule", base.CreateModule)
		group.DELETE("/allmodule", base.DeleteGroupAndModule)
	}

	// module_entry
	{
		group.GET("module_entry", base.QueryModuleEntry)
		group.POST("module_entry", base.CreateModuleEntry)
		group.DELETE("module_entry", base.DeleteModuleEntry)
	}

	//recent_visit
	{
		group.GET("recent_visit", base.ListRecentVisit)
	}

	//group.Use(recentvisit.RecentVisit(server))
	//showallgroupmodule
	{
		group.GET("showallgroup", base.ListGroup)
		group.GET("showallmodule", base.ListModule)
	}

	return base
}
