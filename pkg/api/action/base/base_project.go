package base

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource"
	"net/http"
	"strconv"
)

func (b *baseServer) ListProject(g *gin.Context) {
	pageInt, _ := strconv.Atoi(g.DefaultQuery("page", "1"))
	pageSizeInt, _ := strconv.Atoi(g.DefaultQuery("pagesize", "10"))

	results, count, err := b.UserProjectService.List(int64(pageInt), int64(pageSizeInt))
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}
	g.JSON(http.StatusOK, map[string]interface{}{"data": results, "count": count})
}

func (b *baseServer) CreateProject(g *gin.Context) {
	rawData, err := g.GetRawData()
	if err != nil {
		api.RequestParamsError(g, "get rawData error", err)
		return
	}
	request := &apiResource.RequestUserProject{}
	if err := json.Unmarshal(rawData, request); err != nil {
		api.RequestParamsError(g, "unmarshal json error", err)
		return
	}

	err = b.UserProjectService.Create(request)
	if err != nil {
		api.RequestParamsError(g, "create user error", err)
		return
	}
	g.JSON(http.StatusOK, request)
}
