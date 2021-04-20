package workorder

import (
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/service"
)

type WorkOrderService struct {
	service.IService
}

func NewWorkOrderService(i service.IService) *WorkOrderService {
	return &WorkOrderService{i}
}

func (s *WorkOrderService) List(orderType int, orderStatus int, search string, page, pageSize int64) ([]interface{}, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{
		"spec.order_type": orderType,
		"spec.order_status": orderStatus,
	}

	sort := map[string]interface{}{
		"spec.update_time": -1,
		"metadata.created_time": -1,
	}

	return s.IService.ListByFilter(common.DefaultNamespace, common.WorkOrder, filter, sort, offset, pageSize)
}

