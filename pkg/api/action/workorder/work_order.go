package workorder

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
	serviceUser "github.com/yametech/devops/pkg/service/base"
)

type WorkOrder struct {
	*api.Server
	*serviceUser.UserService
}

func NewWorkOrder(serviceName string, server *api.Server) *WorkOrder {

	workOrder := &WorkOrder{
		Server:      server,
		UserService: serviceUser.NewUser(server.IService),
	}
	group := workOrder.Group(fmt.Sprintf("/%s", serviceName))

	{
		group.GET("/workorder", workOrder.ListWorkOrder)
	}

	return workOrder
}
