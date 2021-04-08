package workorder

import (
	"fmt"
	"github.com/yametech/devops-zpk-server/pkg/api"
	serviceUser "github.com/yametech/devops-zpk-server/pkg/service/user"
)

type WorkOrder struct {
	*api.Server
	*serviceUser.User
}

func NewWorkOrder(serviceName string, server *api.Server) *WorkOrder {

	workOrder := &WorkOrder{
		Server: server,
		User:   serviceUser.NewUser(server.IService),
	}
	group := workOrder.Group(fmt.Sprintf("/%s", serviceName))

	{
		group.GET("/workorder", workOrder.ListWorkOrder)
	}

	return workOrder
}
