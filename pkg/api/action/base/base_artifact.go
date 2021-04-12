package base

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	"net/http"
)

func (b *baseServer) ListArtifact(g *gin.Context) {
	data, err := b.User.ListAr()
	if err != nil {
		api.RequestParamsError(g, "delete fail", err)
		return
	}
	g.JSON(http.StatusOK, data)
}
