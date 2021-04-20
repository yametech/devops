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
		group.GET("/:name", func(c *gin.Context) {
			name := c.Param("name")
			message := fmt.Sprintf("hello world %s ", name)
			c.JSON(200, gin.H{
				"message": message,
				"code":    200,
			})
		})
	}
	return artifactory
}
