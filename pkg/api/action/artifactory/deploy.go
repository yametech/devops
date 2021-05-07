package artifactory

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/artifactory"
	"io"
	"strconv"
)

func (b *baseServer) WatchDeploy(g *gin.Context) {
	version := g.DefaultQuery("version", "0")
	objectChan, closed := b.DeployService.Watch(version)

	streamEndEvent := "STREAM_END"
	g.Stream(func(w io.Writer) bool {
		select {
		case <-g.Writer.CloseNotify():
			closed <- struct{}{}
			close(closed)
			g.SSEvent("", streamEndEvent)
			return false
		case object, ok := <-objectChan:
			if !ok {
				g.SSEvent("", streamEndEvent)
				return false
			}
			g.SSEvent("", object)
		}
		return true
	},
	)

}

func (b *baseServer) ListDeploy(g *gin.Context) {
	pageInt, _ := strconv.Atoi(g.DefaultQuery("page", "1"))
	pageSizeInt, _ := strconv.Atoi(g.DefaultQuery("pagesize", "10"))
	appName := g.DefaultQuery("appname", "")

	results, count, err := b.DeployService.List(appName, int64(pageInt), int64(pageSizeInt))
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	data := map[string]interface{}{"results": results, "count": count}
	api.ResponseSuccess(g, data, "")
}

func (b *baseServer) CreateDeploy(g *gin.Context) {
	request := &apiResource.RequestDeploy{}
	err := g.ShouldBindJSON(request)
	if err != nil {
		api.RequestParamsError(g, "bind json error", err)
	}

	err = b.DeployService.Create(request)
	if err != nil {
		api.RequestParamsError(g, "create deploy error", err)
		return
	}
	api.ResponseSuccess(g, request, "")

}
