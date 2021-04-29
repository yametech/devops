package controller

import (
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/proc"
	"github.com/yametech/devops/pkg/store"
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

func (a *AppServiceController) recvWorkOrder(errors chan<- error) {
	version := time.Now().Unix()
	workOrderCoder := store.GetResourceCoder(string(common.WorkOrder))
	_ = version
	_ = workOrderCoder
}
