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
	name := g.DefaultQuery("name", "")

	results,  err := b.UserService.List(name, int64(pageInt), int64(pageSizeInt))
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}
	g.JSON(http.StatusOK, map[string]interface{}{"data": results})
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

	err = b.UserService.Create(request)
	if err != nil {
		api.RequestParamsError(g, "create user error", err)
		return
	}
	g.JSON(http.StatusOK, request)
}

func (b *baseServer) GetUser(g *gin.Context) {
	uuid := g.Param("uuid")
	data, err := b.UserService.GetByUUID(uuid)
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}
	g.JSON(http.StatusOK, data)
}

func (b *baseServer) DeleteUser(g *gin.Context) {
	uuid := g.Param("uuid")
	err := b.UserService.Delete(uuid)
	if err != nil {
		api.RequestParamsError(g, "delete fail", err)
		return
	}
	g.JSON(http.StatusOK, nil)
}

func (b *baseServer) UpdateUser(g *gin.Context) {
	uuid := g.Param("uuid")

	rawData, err := g.GetRawData()
	if err != nil {
		api.RequestParamsError(g, "get rawData error", err)
		return
	}
	request := &apiResource.RequestUser{}
	if err := json.Unmarshal(rawData, &request); err != nil {
		api.RequestParamsError(g, "unmarshal json error", err)
		return
	}

	user, _, err := b.UserService.Update(uuid, request)
	if err != nil {
		api.RequestParamsError(g, "update fail", err)
		return
	}

	g.JSON(http.StatusOK, user)

}
