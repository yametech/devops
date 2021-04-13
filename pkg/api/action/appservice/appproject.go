package appservice

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) ListAppProject(g *gin.Context) {
	g.JSON(http.StatusOK, "")
}
