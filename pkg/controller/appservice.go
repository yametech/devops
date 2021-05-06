package controller

import (
	"fmt"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/proc"
	"github.com/yametech/devops/pkg/resource/workorder"
	"github.com/yametech/devops/pkg/store"
	"github.com/yametech/devops/pkg/utils"
	"log"
	"time"
)

var _ Controller = &AppServiceController{}

type AppServiceController struct {
	store.IKVStore
	proc *proc.Proc
}

func NewPipelineController(store store.IKVStore) *AppServiceController {
	server := &AppServiceController{
		IKVStore: store,
		proc:     proc.NewProc(),
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

}
