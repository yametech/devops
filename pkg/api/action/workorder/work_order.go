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
		group.GET("/workorder", workOrder.ListWorkOrder)
		group.POST("/workorder", workOrder.CreateWorkOrder)
		group.GET("/workorder/:uuid", workOrder.GetWorkOrder)
		group.PUT("/workorder/:uuid", workOrder.UpdateWorkOrder)
		group.DELETE("/workorder/:uuid", workOrder.DeteleWorkOrder)
	}

	return workOrder
}
