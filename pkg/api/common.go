package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequestParamsError(g *gin.Context, message string, err error) {
	g.JSON(http.StatusBadRequest, gin.H{"message": message, "error": err.Error()})
	g.Abort()
}


func ResponseError(g *gin.Context, err error) {
	g.JSON(http.StatusOK, gin.H{"msg": err.Error(), "code": http.StatusBadRequest, "data": nil})
	g.Abort()
}

func ResponseSuccess(g *gin.Context, data interface{}){
	g.JSON(http.StatusOK, gin.H{"msg": "", "code": http.StatusOK, "data": data})
	g.Abort()
}
