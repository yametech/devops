package base

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource"
	"net/http"
)

func (b *baseServer) ListUser(g *gin.Context) {
	data, err := b.User.Query("users", nil)
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}
	g.JSON(http.StatusOK, data)
}

func (b *baseServer) GetUser(g *gin.Context) {
	name, _ := g.GetQuery("name")
	user := &resource.User{}
	err := b.User.QueryOne("users", map[string]interface{}{"name": name}, user)
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
	user := &resource.User{
		Metadata: core.Metadata{
			Name: request.Name,
			Kind: request.Kind,
		},
		Spec: resource.UserSpec{
			NickName: request.NickName,
			Username: request.Username,
			Password: request.Password,
		},
	}
	user.GenerateVersion()

	err = b.User.Create("", user)
	if err != nil {
		api.RequestParamsError(g, "get data error", err)
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
	request := &apiResource.RequestUser{}
	if err := json.Unmarshal(rawData, request); err != nil {
		api.RequestParamsError(g, "unmarshal json error", err)
		return
	}
	err = b.User.Update("users", uuid, request)
	if err != nil {
		api.RequestParamsError(g, "update fail", err)
		return
	}
	g.JSON(http.StatusOK, request)

}

func (b *baseServer) DeleteUser(g *gin.Context) {
	uuid := g.Param("uuid")
	user := &resource.User{}
	err := b.User.QueryOne("users", map[string]interface{}{"uuid": uuid}, user)
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}
	err = b.User.Delete("users", user)
	if err != nil {
		api.RequestParamsError(g, "delete fail", err)
		return
	}
	g.JSON(http.StatusOK, nil)
}
