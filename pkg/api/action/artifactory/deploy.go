package artifactory

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/artifactory"
	"strconv"
)

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
