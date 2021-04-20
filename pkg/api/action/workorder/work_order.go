package workorder

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/service/workorder"
)

type WorkOrder struct {
	*api.Server
	*workorder.WorkOrderService
}

func NewWorkOrder(serviceName string, server *api.Server) *WorkOrder {

	workOrder := &WorkOrder{
		Server:      server,
		WorkOrderService: workorder.NewWorkOrderService(server.IService),
	}
	group := workOrder.Group(fmt.Sprintf("/%s", serviceName))

	{
		group.GET("/workorder", workOrder.ListWorkOrder)
		group.POST("/workorder", workOrder.CreateWorkOrder)
	}

	return workOrder
}
