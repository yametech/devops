package base

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource"
	"net/http"
)

func (b *baseServer) ListUser(g *gin.Context) {
	data, err := b.User.List()
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}
	g.JSON(http.StatusOK, data)
}

func (b *baseServer) GetUser(g *gin.Context) {
	name, _ := g.GetQuery("name")
	user, err := b.User.Query(map[string]interface{}{"name": name})
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}

	g.JSON(http.StatusOK, user)
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
		api.RequestParamsError(g, "create user error", err)
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
