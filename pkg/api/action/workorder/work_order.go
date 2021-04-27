package workorder

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/service/workorder"
)

type WorkOrder struct {
	*api.Server
	*workorder.Service
}

func NewWorkOrder(serviceName string, server *api.Server) *WorkOrder {

	workOrder := &WorkOrder{
		Server:  server,
		Service: workorder.NewWorkOrderService(server.IService),
	}
	group := workOrder.Group(fmt.Sprintf("/%s", serviceName))

	{
		group.GET("/order", workOrder.ListWorkOrder)
		group.POST("/order", workOrder.CreateWorkOrder)
		group.GET("/order/:uuid", workOrder.GetWorkOrder)
		group.PUT("/order/:uuid", workOrder.UpdateWorkOrder)
		group.DELETE("/order/:uuid", workOrder.DeteleWorkOrder)
		group.GET("/status", workOrder.GetWorkOrderStatus)
	}

	return workOrder
}
