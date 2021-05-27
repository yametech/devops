package artifactory

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/artifactory"
	"github.com/yametech/devops/pkg/resource/artifactory"
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

	query := map[string]interface{}{}
	query["spec.app_name"] = g.DefaultQuery("appname", "")
	query["spec.deploy_namespace"] = g.DefaultQuery("namespace", "")
	query["spec.create_team"] = g.DefaultQuery("team", "")
	query["spec.create_user_id"] = g.DefaultQuery("username", "")

	results, count, err := b.DeployService.List(query, int64(pageInt), int64(pageSizeInt))
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
		return
	}

	//todo:待网关能传中文名字后，要获取中文名
	request.UserName = g.Request.Header.Get("X-Wrapper-Username")

	err = b.DeployService.Create(request)
	if err != nil {
		api.RequestParamsError(g, "create deploy error", err)
		return
	}
	api.ResponseSuccess(g, request, "")
}

func (b *baseServer) GetDeploy(g *gin.Context) {
	uuid := g.Param("uuid")
	ns := g.DefaultQuery("namespace", "")
	//latest := g.DefaultQuery("latest", "")

	var err error
	data := &artifactory.Deploy{}

	if ns != "" {
		data, err = b.DeployService.GetByAppName(uuid, ns)
	} else {
		data, err = b.DeployService.GetByUUID(uuid)
	}

	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}

	result := map[string]interface{}{"results": data}
	api.ResponseSuccess(g, result, "")
}
