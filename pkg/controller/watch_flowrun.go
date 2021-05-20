package controller

import (
	"encoding/json"
	"fmt"
	"github.com/r3labs/sse/v2"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/proc"
	"github.com/yametech/devops/pkg/resource/artifactory"
	"github.com/yametech/devops/pkg/store"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var _ Controller = &WatchFlowRun{}

type WatchFlowRun struct {
	store.IKVStore
	proc *proc.Proc
}

func NewWatchFlowRun(ikvStore store.IKVStore) *WatchFlowRun {
	server := &WatchFlowRun{
		IKVStore: ikvStore,
		proc:     proc.NewProc(),
	}
	return server
}

func (w *WatchFlowRun) Run() error {
	fmt.Println(fmt.Sprintf("[Controller]%v start --> %v", reflect.TypeOf(w), time.Now()))
	w.proc.Add(w.artifactConnect)
	return <-w.proc.Start()
}

func (w *WatchFlowRun) artifactConnect(errC chan<- error) {
	//GetOldFlownRun
	fmt.Printf("[Controller]%v begin handle old flowrun.\n", reflect.TypeOf(w))

	resp, err := http.Get(fmt.Sprintf("%s/flowrun/", common.EchoerUrl))
	if err != nil {
		fmt.Printf("[Controller]%v handle old flowrun get url fail, err %s.\n", reflect.TypeOf(w), err)
		errC <- err
	}
	if resp == nil {
		fmt.Printf("[Controller]%v handle old flowrun res was empty.\n", reflect.TypeOf(w))
	}
	defer resp.Body.Close()
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Printf("[Controller]%v handle old flowrun read resp fail, err %s.\n", reflect.TypeOf(w), err1)
		errC <- err1
	}
	var uBody []*FlowRun
	if err := json.Unmarshal(body, &uBody); err != nil {
		fmt.Printf("[Controller]%v handle old flowrun Unmarshal fail, err %s.\n", reflect.TypeOf(w), err)
		errC <- err
	}
	for _, oneFlowRun := range uBody {
		w.HandleFlowRun(oneFlowRun)
	}

	//artifactConnect
	fmt.Printf("[Controller]%v begin watch echoer.\n", reflect.TypeOf(w))
	Version := time.Now().Unix() - 100
	url := fmt.Sprintf("%s/watch?resource=flowrun?version=%d", common.EchoerUrl, Version)
	fmt.Printf("[Controller]%v begin func artifactConnect and url: %s.\n", reflect.TypeOf(w), url)

	client := sse.NewClient(url)
	err = client.SubscribeRaw(func(msg *sse.Event) {
		a := &FlowRun{}
		err := json.Unmarshal(msg.Data, &a)
		if err != nil {
			fmt.Println("Error When artifactConnect Unmarshal:", err)
			errC <- err
		}
		fmt.Printf("[Controller]%v watching echoer: %v\n", reflect.TypeOf(w), a.Metadata)
		go w.HandleFlowRun(a)
		//b, _ := strconv.Atoi(a["metadata"]["version"])
		//fmt.Println(b)
	})
	if err != nil {
		errC <- err
	}
}

func (w *WatchFlowRun) HandleFlowRun(run *FlowRun) {
	flowRunName := run.Metadata.Name
	//Determine the type of split
	if strings.Contains(flowRunName, "_") {
		fmt.Printf("[Controller]%v HandleFlowrun get flowRun %s \n", reflect.TypeOf(w), run.Metadata.Name)
		keyValue := strings.Split(flowRunName, "_")
		if len(keyValue) != 2 || keyValue[0] != common.DefaultNamespace {
			fmt.Printf("[Controller]%v HandleFlowrun result of split isn't two or not belongs to devops %s \n", reflect.TypeOf(w), run.Metadata.Name)
			return
		}
		for _, flowStep := range run.Spec.Steps {
			//action is Done
			if flowStep.Spec.ActionRun.Done != true {
				continue
			}
			temp := strings.Split(flowStep.Metadata.Name, "_")
			stepUUID := temp[1]
			actionType := temp[0]
			//switch actionType from actionName
			switch actionType {
			//actionName EchoerCI
			case common.CI:
				step := &artifactory.Artifact{}
				err := w.GetByUUID(common.DefaultNamespace, common.Artifactory, stepUUID, step)
				if err != nil {
					fmt.Printf("[Controller]%v error: action: %s, uuid: %s ,err: %s \n", reflect.TypeOf(w), flowStep.Spec.ActionRun.ActionName, stepUUID, err.Error())
					continue
				}
				if flowStep.Spec.Response.State == "SUCCESS" {
					step.Spec.ArtifactStatus = artifactory.Built
				} else if flowStep.Spec.Response.State == "FAIL" {
					step.Spec.ArtifactStatus = artifactory.BuiltFAIL
				} else if flowStep.Spec.Response.State == "TIMEOUT" {
					step.Spec.ArtifactStatus = artifactory.BuildTimeOut
				}
				_, _, err = w.Apply(common.DefaultNamespace, common.Artifactory, step.Metadata.UUID, step, false)
				if err != nil {
					fmt.Printf("[Controller]%v error: action: %s, uuid: %s ,err: %s \n", reflect.TypeOf(w), flowStep.Spec.ActionRun.ActionName, stepUUID, err.Error())
				}
				fmt.Printf("[Controller]%v msg: action: %s, uuid: %s status updata done! \n", reflect.TypeOf(w), flowStep.Spec.ActionRun.ActionName, stepUUID)
			case common.CD:
				step := &artifactory.Deploy{}
				err := w.GetByUUID(common.DefaultNamespace, common.Deploy, stepUUID, step)
				if err != nil {
					fmt.Printf("[Controller]%v error: action: %s, uuid: %s ,err: %s \n", reflect.TypeOf(w), flowStep.Spec.ActionRun.ActionName, stepUUID, err.Error())
					continue
				}
				if flowStep.Spec.Response.State == "SUCCESS" {
					step.Spec.DeployStatus = artifactory.Deployed
				} else if flowStep.Spec.Response.State == "FAIL" {
					step.Spec.DeployStatus = artifactory.DeployFail
				} else if flowStep.Spec.Response.State == "TIMEOUT" {
					step.Spec.DeployStatus = artifactory.DeployTimeOut
				}
				_, _, err = w.Apply(common.DefaultNamespace, common.Deploy, step.Metadata.UUID, step, false)
				if err != nil {
					fmt.Printf("[Controller]%v error: action: %s, uuid: %s ,err: %s \n", reflect.TypeOf(w), flowStep.Spec.ActionRun.ActionName, stepUUID, err.Error())
				}
				fmt.Printf("[Controller]%v msg: action: %s, uuid: %s status updata done! \n", reflect.TypeOf(w), flowStep.Spec.ActionRun.ActionName, stepUUID)
			default:
				fmt.Printf("[Controller]%v msg: actiontype is not found: %s, uuid: %s \n", reflect.TypeOf(w), actionType, stepUUID)
			}
		}
	}
}
