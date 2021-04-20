package controller

import (
	"encoding/json"
	"fmt"
	"github.com/r3labs/sse/v2"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/resource/artifactory"
	"github.com/yametech/devops/pkg/store"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var _ Controller = &WatchFlowRun{}

type FlowRun struct {
	Metadata struct {
		Name    string      `json:"name"`
		Kind    string      `json:"kind"`
		Version int         `json:"version"`
		UUID    string      `json:"uuid"`
		Labels  interface{} `json:"labels"`
	} `json:"metadata"`
	Spec struct {
		Steps []struct {
			Metadata struct {
				Name    string      `json:"name"`
				Kind    string      `json:"kind"`
				Version int         `json:"version"`
				UUID    string      `json:"uuid"`
				Labels  interface{} `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				FlowID      string `json:"flow_id"`
				FlowRunUUID string `json:"flow_run_uuid"`
				ActionRun   struct {
					ActionName   string `json:"action_name"`
					ActionParams struct {
						Branch      string `json:"branch"`
						CodeType    string `json:"codeType"`
						CommitID    string `json:"commitId"`
						GitURL      string `json:"gitUrl"`
						Output      string `json:"output"`
						ProjectFile string `json:"projectFile"`
						ProjectPath string `json:"projectPath"`
						RetryCount  int    `json:"retryCount"`
						ServiceName string `json:"serviceName"`
					} `json:"action_params"`
					ReturnStateMap struct {
						FAIL    string `json:"FAIL"`
						SUCCESS string `json:"SUCCESS"`
					} `json:"return_state_map"`
					Done bool `json:"done"`
				} `json:"action_run"`
				Response struct {
					State string `json:"state"`
				} `json:"response"`
				Data            string      `json:"data"`
				RetryCount      int         `json:"retry_count"`
				GlobalVariables interface{} `json:"global_variables"`
			} `json:"spec"`
		} `json:"steps"`
		HistoryStates  []string    `json:"history_states"`
		LastState      string      `json:"last_state"`
		CurrentState   string      `json:"current_state"`
		LastEvent      string      `json:"last_event"`
		LastErr        string      `json:"last_err"`
		GlobalVariable interface{} `json:"global_variable"`
	} `json:"spec"`
}

type WatchFlowRun struct {
	store.IKVStore
}

func NewWatchFlowRun(ikvStore store.IKVStore) *WatchFlowRun {
	server := &WatchFlowRun{
		ikvStore,
	}
	return server
}

func (w *WatchFlowRun) Run() error {
	fmt.Println(fmt.Sprintf("[Controller]%v start --> %v", reflect.TypeOf(w), time.Now()))
	errC := make(chan error)
	go w.GetOldVersion(errC)
	go w.ArtifactConnect(errC)
	return <-errC
}

func (w *WatchFlowRun) GetOldVersion(errC chan<- error) {
	fmt.Printf("[Controller]%v begin func GetOldVersion.\n", reflect.TypeOf(w))
	resp, err := http.Get(fmt.Sprintf("%s/flowrun/", common.EchoerUrl))
	if err != nil {
		errC <- err
	}
	if resp == nil {
		fmt.Printf("[Controller]%v GetOldVersion res was empty.\n", reflect.TypeOf(w))
		return
	}
	defer resp.Body.Close()
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		errC <- err1
	}
	var uBody []FlowRun
	if err := json.Unmarshal(body, &uBody); err != nil {
		errC <- err
	}
	//fmt.Println(string(body))
	//for _, a := range uBody {
	//	b, _ := json.Marshal(a)
	//	fmt.Println(string(b))
	//}
	//fmt.Println(uBody[len(uBody)-1])
	for _, oneFlowRun := range uBody {
		go w.HandleFlowrun(oneFlowRun)
	}
}

func (w *WatchFlowRun) ArtifactConnect(errC chan<- error) {
	Version := time.Now().Unix()
	url := fmt.Sprintf("%s/watch?resource=flowrun?version=%d", common.EchoerUrl, Version)
	fmt.Printf("[Controller]%v begin func ArtifactConnect and url: %s.\n", reflect.TypeOf(w), url)
	client := sse.NewClient(url)
	err := client.SubscribeRaw(func(msg *sse.Event) {
		var a FlowRun
		err := json.Unmarshal(msg.Data, &a)
		//err = errors.New("new error")
		if err != nil {
			fmt.Println("Error When ArtifactConnect Unmarshal:", err)
			errC <- err
		}
		fmt.Printf("%v\n", a.Metadata)
		w.HandleFlowrun(a)
		//b, _ := strconv.Atoi(a["metadata"]["version"])
		//fmt.Println(b)
	})
	if err != nil {
		errC <- err
	}
}

func (w *WatchFlowRun) HandleFlowrun(run FlowRun) {
	flowRunName := run.Metadata.Name
	fmt.Printf("[Controller]%v HandleFlowrun get flowRun %s \n", reflect.TypeOf(w), run.Metadata.Name)
	//Determine the type of split
	switch {
	case strings.Contains(flowRunName, "_"):
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
			stepUUID := strings.Split(flowStep.Metadata.Name, "_")[1]
			//switch actionType from actionName
			switch flowStep.Spec.ActionRun.ActionName {
			//actionName EchoerCI
			case common.EchoerCI:
				step := &artifactory.Artifact{}
				err := w.GetByUUID(common.DefaultNamespace, common.Artifactory, stepUUID, step)
				if err != nil {
					fmt.Printf("[Controller]%v error: step: %s, uuid: %s ,err: %s \n", reflect.TypeOf(w), common.EchoerCI, stepUUID, err.Error())
					continue
				}
				if flowStep.Spec.Response.State == "SUCCESS" && step.Spec.ArtifactStatus != artifactory.Built {
					step.Spec.ArtifactStatus = artifactory.Built
				} else if flowStep.Spec.Response.State == "FAIL" && step.Spec.ArtifactStatus != artifactory.BuiltFAIL {
					step.Spec.ArtifactStatus = artifactory.BuiltFAIL
				}
				_, _, err = w.Apply(common.DefaultNamespace, common.Artifactory, step.Metadata.UUID, step)
				if err != nil {
					fmt.Printf("[Controller]%v error: step: %s, uuid: %s ,err: %s \n", reflect.TypeOf(w), common.EchoerCI, stepUUID, err.Error())
				}
			//TODO:actionName EchoerCI
			case common.EchoerCD:
				fmt.Println("TODO CD")
			}
		}
	}
}
