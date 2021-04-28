package artifactory

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/artifactory"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (b *baseServer) WatchAr(g *gin.Context) {
	version := g.DefaultQuery("version", "0")
	objectChan, closed := b.ArtifactService.Watch(version)

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

func (b *baseServer) ListArtifact(g *gin.Context) {
	pageInt, _ := strconv.Atoi(g.DefaultQuery("page", "1"))
	pageSizeInt, _ := strconv.Atoi(g.DefaultQuery("pagesize", "10"))
	name := g.DefaultQuery("name", "")

	results, count, err := b.ArtifactService.List(name, int64(pageInt), int64(pageSizeInt))
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}

	data := map[string]interface{}{"results": results}
	data["count"] = count
	api.ResponseSuccess(g, data, "")
}

func (b *baseServer) CreateArtifact(g *gin.Context) {
	rawData, err := g.GetRawData()
	if err != nil {
		api.RequestParamsError(g, "get rawData error", err)
		return
	}
	request := &apiResource.RequestArtifact{}
	if err := json.Unmarshal(rawData, request); err != nil {
		api.RequestParamsError(g, "unmarshal json error", err)
		return
	}

	res, err := b.ArtifactService.Create(request)
	if err != nil {
		api.RequestParamsError(g, "create artifact error", err)
		return
	}
	api.ResponseSuccess(g, res, "")
}

func (b *baseServer) GetArtifact(g *gin.Context) {
	uuid := g.Param("uuid")
	data, err := b.ArtifactService.GetByUUID(uuid)
	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}

	result := map[string]interface{}{"results": data}
	api.ResponseSuccess(g, result, "")
}

func (b *baseServer) DeleteArtifact(g *gin.Context) {
	uuid := g.Param("uuid")
	err := b.ArtifactService.Delete(uuid)
	if err != nil {
		api.RequestParamsError(g, "delete fail", err)
		return
	}
	api.ResponseSuccess(g, nil, "删除成功")
}

func (b *baseServer) UpdateArtifact(g *gin.Context) {
	uuid := g.Param("uuid")

	rawData, err := g.GetRawData()
	if err != nil {
		api.RequestParamsError(g, "get rawData error", err)
		return
	}
	request := &apiResource.RequestArtifact{}
	if err := json.Unmarshal(rawData, &request); err != nil {
		api.RequestParamsError(g, "unmarshal json error", err)
		return
	}

	user, _, err := b.ArtifactService.Update(uuid, request)
	if err != nil {
		api.RequestParamsError(g, "update fail", err)
		return
	}

	g.JSON(http.StatusOK, user)

}

func (b *baseServer) GetBranchList(g *gin.Context) {
	gitPath := g.Query("gitpath")
	gitPath = strings.Replace(gitPath, ".git", "", -1)
	sliceTemp := strings.Split(gitPath, "/")
	org, name := "", ""
	if len(sliceTemp) >= 2 {
		org = sliceTemp[len(sliceTemp)-2]
		name = sliceTemp[len(sliceTemp)-1]
	} else {
		return
	}

	results, err := b.ArtifactService.GetBranch(gitPath)

	if err != nil {
		api.RequestParamsError(g, "error", err)
		return
	}

	data := map[string]interface{}{"results": results}
	data["count"] = len(results)
	api.ResponseSuccess(g, data, "")
}
