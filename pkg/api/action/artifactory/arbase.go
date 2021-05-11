package artifactory

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	arbaseService "github.com/yametech/devops/pkg/service/artifactory"
)

type baseServer struct {
	*api.Server
	*arbaseService.ArtifactService
	*arbaseService.DeployService
}

func NewArBaseServer(serviceName string, server *api.Server) *baseServer {
	base := &baseServer{
		Server:          server,
		ArtifactService: arbaseService.NewArtifact(server.IService),
		DeployService:   arbaseService.NewDeployService(server.IService),
	}
	group := base.Group(fmt.Sprintf("/%s", serviceName))

	// watch
	{
		group.GET("/artifactwatch", base.WatchAr)
		group.GET("/deploywatch", base.WatchDeploy)
	}

	//UserProjectService
	{
		group.GET("/artifact", base.ListArtifact)
		group.GET("/artifact/:uuid", base.GetArtifact)
		group.POST("/artifact", base.CreateArtifact)
		group.PUT("/artifact/:uuid", base.UpdateArtifact)
		group.DELETE("/artifact/:uuid", base.DeleteArtifact)
	}

	//GetBranch
	{
		group.GET("/getbranch", base.GetBranchList)
	}

	//Deploy
	{
		group.GET("/deploy", base.ListDeploy)
		group.POST("/deploy", base.CreateDeploy)

	}
	return base
}
