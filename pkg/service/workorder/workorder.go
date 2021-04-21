package workorder

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/workorder"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/workorder"
	"github.com/yametech/devops/pkg/service"
)

type Service struct {
	service.IService
}

func NewWorkOrderService(i service.IService) *Service {
	return &Service{i}
}

func (s *Service) List(orderType int, search string, page, pageSize int64) ([]interface{}, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{
		"spec.order_type": orderType,
		//"spec.order_status": orderStatus,
	}

	sort := map[string]interface{}{
		"metadata.version":      -1,
		"metadata.created_time": -1,
	}

	return s.IService.ListByFilter(common.DefaultNamespace, common.WorkOrder, filter, sort, offset, pageSize)
}

func (s *Service) Create(request *apiResource.Request) (core.IObject, error) {
	req := &workorder.WorkOrder{
		Spec: workorder.Spec{
			OrderType: request.OrderType,
			Title:     request.Title,
			Attribute: request.Attribute,
			Apply:     request.Apply,
			Check:     request.Check,
			Result:    request.Result,
		},
	}

	req.GenerateNumber()
	req.GenerateVersion()
	return s.IService.Create(common.DefaultNamespace, common.WorkOrder, req)
}

func (s *Service) Get(uuid string) (core.IObject, error) {
	order := &workorder.WorkOrder{}
	if err := s.IService.GetByUUID(common.DefaultNamespace, common.WorkOrder, uuid, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) Update(uuid string, request *apiResource.Request) (core.IObject, bool, error) {
	dbObj := &workorder.WorkOrder{}
	if err := s.GetByUUID(common.DefaultNamespace, common.WorkOrder, uuid, dbObj); err != nil {
		return nil, false, errors.New("The workorder is not exist")
	}
	dbObj.Spec.OrderType = request.OrderType
	dbObj.Spec.Title = request.Title
	dbObj.Spec.Attribute = request.Attribute
	dbObj.Spec.Apply = request.Apply
	dbObj.Spec.Check = request.Check
	dbObj.Spec.Result = request.Result

	dbObj.GenerateVersion()
	return s.IService.Apply(common.DefaultNamespace, common.WorkOrder, dbObj.UUID, dbObj, false)
}

func (s *Service) Delete(uuid string) (bool, error) {
	if err := s.IService.Delete(common.DefaultNamespace, common.WorkOrder, uuid); err != nil {
		return false, err
	}

	return true, nil
}
