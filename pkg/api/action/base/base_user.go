package base

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (b *baseServer) CreateUser(g *gin.Context) {
	g.JSON(http.StatusOK, "")
}
