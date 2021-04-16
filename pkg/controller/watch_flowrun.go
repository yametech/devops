package controller

import (
	"encoding/json"
	"fmt"
	"github.com/r3labs/sse/v2"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/store"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

var Version int

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

func (w WatchFlowRun) Run() error {
	fmt.Println(fmt.Sprintf("[Controller]%v start --> %v", reflect.TypeOf(w), time.Now()))
	errC := make(chan error)
	w.FirstConnect(errC)
	w.ArtifactConnect(errC)
	return <-errC
}

func (w WatchFlowRun) FirstConnect(errC chan<- error) {
	resp, err := http.Get(fmt.Sprintf("%s/flowrun/", common.EchoerUrl))
	if err != nil {
		errC <- err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errC <- err
	}
	var uBody []FlowRun
	//fmt.Println(string(body))
	if err := json.Unmarshal(body, &uBody); err != nil {
		errC <- err
	}
	fmt.Println(uBody)
	Version = uBody[0].Metadata.Version - 100
}

func (w WatchFlowRun) ArtifactConnect(errC chan<- error) {
	client := sse.NewClient(fmt.Sprintf("%s/watch?resource=flowrun?version=%d", common.EchoerUrl, Version))
	err := client.SubscribeRaw(func(msg *sse.Event) {
		var a FlowRun
		err := json.Unmarshal(msg.Data, &a)
		if err != nil {
			fmt.Println("Error When ArtifactConnect Unmarshal:", err)
			errC <- err
		}
		fmt.Printf("%v\n", a.Metadata)
		//b, _ := strconv.Atoi(a["metadata"]["version"])
		//fmt.Println(b)
	})
	if err != nil {
		errC <- err
	}
}

func HandleFlowrun() {

}
