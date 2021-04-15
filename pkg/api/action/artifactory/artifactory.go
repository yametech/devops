package artifactory

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
)

type artifactoryServer struct {
	*api.Server
}

func NewArtifactoryServer(serviceName string, server *api.Server) *artifactoryServer {
	artifactory := &artifactoryServer{
		server,
	}
	group := artifactory.Group(fmt.Sprintf("/%s", serviceName))

	//ArtifactoryServer
	{
		group.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "hello world",
				"code":    200,
			})
		})
	}
	return artifactory
}
