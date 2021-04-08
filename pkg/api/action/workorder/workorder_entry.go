package workorder

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (w *WorkOrder) ListWorkOrder(g *gin.Context) {
	g.JSON(http.StatusOK, "")
}
