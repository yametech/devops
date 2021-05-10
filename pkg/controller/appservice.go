package controller

import (
	"fmt"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/proc"
	"github.com/yametech/devops/pkg/resource/workorder"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/service/appservice"
	"github.com/yametech/devops/pkg/store"
	"github.com/yametech/devops/pkg/utils"
	"log"
	"time"
)

var _ Controller = &AppServiceController{}

type AppServiceController struct {
	store.IKVStore
	proc *proc.Proc
	handlerMap map[workorder.OrderType]map[workorder.OrderStatus]func(obj *workorder.WorkOrder) error
}

func NewPipelineController(store store.IKVStore) *AppServiceController {
	baseService := service.NewBaseService(store)
 	appConfigService := appservice.NewAppConfigService(baseService)
 	namespaceService := appservice.NewNamespaceService(baseService)
	rsHandlerMap := map[workorder.OrderStatus]func(obj *workorder.WorkOrder) error{
		workorder.Checking: appConfigService.OrderToResourceCheck,
		workorder.Rejected: appConfigService.OrderToResourceFailed,
		workorder.Finish: appConfigService.OrderToResourceSuccess,
	}

	nsHandlerMap := map[workorder.OrderStatus]func(obj *workorder.WorkOrder) error{
		workorder.Finish: namespaceService.OrderToNamespaceSuccess,
	}

	server := &AppServiceController{
		IKVStore: store,
		proc:     proc.NewProc(),
		handlerMap: map[workorder.OrderType]map[workorder.OrderStatus]func(obj *workorder.WorkOrder) error{
			workorder.Resources: rsHandlerMap,
			workorder.Namespace: nsHandlerMap,
		},
	}


	return server
}

func (a *AppServiceController) Run() error {
	a.proc.Add(a.recvWorkOrder)

	return <-a.proc.Start()
}

func (a *AppServiceController) recvWorkOrder(errC chan<- error) {

	version := time.Now().Unix()
	workOrderCoder := store.GetResourceCoder(string(workorder.WorkerOrderKind))
	if workOrderCoder == nil {
		errC <- fmt.Errorf("(%s) %s", workorder.WorkerOrderKind, "coder not exist")
	}
	workOrderWatchChan := store.NewWatch(workOrderCoder)
	a.Watch2(common.DefaultNamespace, common.WorkOrder, version, workOrderWatchChan)
	log.Println("workOrderController start watching workOrder")

	for {
		select {
		case item, ok := <-workOrderWatchChan.ResultChan():
			if !ok {
				errC <- fmt.Errorf("recvPipeLine watch channal close")
			}
			if item.GetUUID() == "" {
				continue
			}
			workOrderObj := &workorder.WorkOrder{}
			if err := utils.UnstructuredObjectToInstanceObj(&item, workOrderObj); err != nil {
				log.Printf("receive pipeline UnmarshalInterfaceToResource error %s\n", err)
				continue
			}
			go a.handleWorkOrder(workOrderObj)
		}
	}

}

func (a *AppServiceController) handleWorkOrder(obj *workorder.WorkOrder) {
	//TODO: get workOrder if its AppService config and apply it

	if order, exist := a.handlerMap[obj.Spec.OrderType]; exist{
		if handler, ok := order[obj.Spec.OrderStatus]; ok {
			if err := handler(obj); err != nil {
				log.Printf("controller handleWorkOrder error: %s\n", err)
			}
		}
	}
}

func (a *AppServiceController) GetCMDBAppService(errC chan<- error)  {
	go func() {

	}()
}