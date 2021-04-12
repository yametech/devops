package base

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource"
	"net/http"
	"strconv"
)

func (b *baseServer) ListUser(g *gin.Context) {
	pageInt, _ := strconv.Atoi(g.DefaultQuery("page", "1"))
	pageSizeInt, _ := strconv.Atoi(g.DefaultQuery("pagesize", "10"))

	data, count, err := b.User.List(pageInt, pageSizeInt)
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}
	g.JSON(http.StatusOK, map[string]interface{}{"data": data, "count": count})
}

func (b *baseServer) GetUser(g *gin.Context) {
	page := g.DefaultQuery("page", "1")
	pageSize := g.DefaultQuery("pagesize", "10")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	name, _ := g.GetQuery("name")
	data, count, err := b.User.Query(map[string]interface{}{"name": name}, pageInt, pageSizeInt)
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}

	g.JSON(http.StatusOK, map[string]interface{}{"data": data, "count": count})
}

func (b *baseServer) CreateUser(g *gin.Context) {
	rawData, err := g.GetRawData()
	if err != nil {
		api.RequestParamsError(g, "get rawData error", err)
		return
	}
	request := &apiResource.RequestUser{}
	if err := json.Unmarshal(rawData, request); err != nil {
		api.RequestParamsError(g, "unmarshal json error", err)
		return
	}
	user, err := b.User.Create(request)
	if err != nil {
		api.RequestParamsError(g, "create base error", err)
		return
	}

	g.JSON(http.StatusOK, user)

}

func (b *baseServer) UpdateUser(g *gin.Context) {
	uuid := g.Param("uuid")

	rawData, err := g.GetRawData()
	if err != nil {
		api.RequestParamsError(g, "get rawData error", err)
		return
	}
	request := make(map[string]interface{}, 0)
	if err := json.Unmarshal(rawData, &request); err != nil {
		api.RequestParamsError(g, "unmarshal json error", err)
		return
	}
	user, err := b.User.Update(uuid, request)
	if err != nil {
		api.RequestParamsError(g, "update fail", err)
		return
	}

	g.JSON(http.StatusOK, user)

}

func (b *baseServer) DeleteUser(g *gin.Context) {
	uuid := g.Param("uuid")
	err := b.User.Delete(uuid)
	if err != nil {
		api.RequestParamsError(g, "delete fail", err)
		return
	}
	g.JSON(http.StatusOK, nil)
}
