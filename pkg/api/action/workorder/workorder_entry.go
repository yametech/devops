package workorder

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/api"
	apiResource "github.com/yametech/devops/pkg/api/resource/workorder"
	"strconv"
)

func (w *WorkOrder) ListWorkOrder(g *gin.Context) {
	orderType, err := strconv.Atoi(g.DefaultQuery("order_type", "0"))
	if err != nil {
		api.ResponseError(g, errors.New("orderType need int type"))
		return
	}
	orderStatus, err := strconv.Atoi(g.DefaultQuery("order_status", "-1"))
	search := g.Query("search")
	page, err := strconv.ParseInt(g.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		api.ResponseError(g, errors.New("page need int type"))
		return
	}
	pageSize, err := strconv.ParseInt(g.DefaultQuery("page_size", "10"), 10, 64)
	if err != nil {
		api.ResponseError(g, errors.New("pageSize need int type"))
		return
	}

	orders, err := w.Service.List(orderType, orderStatus, search, page, pageSize)

	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, orders, "")
}

func (w *WorkOrder) CreateWorkOrder(g *gin.Context) {
	// Get the user
	user := g.Request.Header.Get("x-wrapper-chinese")

	data := &apiResource.Request{}
	if err := g.ShouldBindJSON(&data); err != nil {
		api.ResponseError(g, err)
		return
	}
	order, err := w.Service.Create(data, user)
	if err != nil {
		api.ResponseError(g, err)
		return
	}
	api.ResponseSuccess(g, order, "")
}

func (w *WorkOrder) GetWorkOrder(g *gin.Context) {
	uuid := g.Param("uuid")
	result, err := w.Service.Get(uuid)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, result, "")
}

func (w *WorkOrder) UpdateWorkOrder(g *gin.Context) {
	uuid := g.Param("uuid")

	data := &apiResource.Request{}
	if err := g.ShouldBindJSON(&data); err != nil {
		api.ResponseError(g, err)
		return
	}

	result, update, err := w.Service.Update(uuid, data)
	if err != nil {
		api.ResponseError(g, err)
		return
	}

	api.ResponseSuccess(g, gin.H{"results": result, "update": update}, "")
}

func (w *WorkOrder) DeleteWorkOrder(g *gin.Context) {
	uuid := g.Param("uuid")

	result, err := w.Service.Delete(uuid)
	if err != nil {
		api.ResponseError(g, err)
		return
	}
	api.ResponseSuccess(g, gin.H{"delete": result}, "")
}
