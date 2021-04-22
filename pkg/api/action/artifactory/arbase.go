package artifactory

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	arbaseService "github.com/yametech/devops/pkg/service/artifactory"
)

type baseServer struct {
	*api.Server
	*arbaseService.ArtifactService
}

func NewArBaseServer(serviceName string, server *api.Server) *baseServer {
	base := &baseServer{
		Server:          server,
		ArtifactService: arbaseService.NewArtifact(server.IService),
	}
	group := base.Group(fmt.Sprintf("/%s", serviceName))

	// watch
	{
		group.GET("/artifactwatch", base.WatchAr)
	}

	//UserProjectService
	{
		group.GET("/artifact", base.ListArtifact)
		group.GET("/artifact/:appname", base.GetArtifact)
		group.POST("/artifact", base.CreateArtifact)
		group.PUT("/artifact/:appname", base.UpdateArtifact)
		group.DELETE("/artifact/:appname", base.DeleteArtifact)
	}

	return base
}
