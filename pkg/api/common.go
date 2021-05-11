package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequestParamsError(g *gin.Context, data string, err error) {
	g.JSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "code": http.StatusBadRequest, "data": data})
	g.Abort()
}

func RequestNotFound(g *gin.Context, data string, err error) {
	g.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "data": data, "msg": err.Error()})
}

func ResponseError(g *gin.Context, err error) {
	g.JSON(http.StatusOK, gin.H{"msg": err.Error(), "code": http.StatusBadRequest, "data": nil})
	g.Abort()
}

func ResponseCodeError(g *gin.Context, err error, code int){
	g.JSON(http.StatusOK, gin.H{"msg": err.Error(), "code": code, "data": nil})
	g.Abort()
}

func ResponseSuccess(g *gin.Context, data interface{}, msg string) {
	g.JSON(http.StatusOK, gin.H{"msg": msg, "code": http.StatusOK, "data": data})
	g.Abort()
}
