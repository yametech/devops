package main

import (
	"encoding/json"
	"fmt"
	"github.com/r3labs/sse"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/store"
	"io/ioutil"
	"net/http"
)

var Version int

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

func (w WatchFlowRun) Run() error {
	panic("implement me")
}

func (w WatchFlowRun) Stop() error {
	panic("implement me")
}

//var _ Controller = &WatchFlowRun{}

func (w WatchFlowRun) firstConnect() error {
	resp, err := http.Get(fmt.Sprintf("%s/flowrun/flow_cd_1618485843706", common.EchoerUrl))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var uBody FlowRun
	//fmt.Println(string(body))
	if err := json.Unmarshal(body, &uBody); err != nil {
		return err
	}
	fmt.Printf("%v", uBody.Metadata)
	Version = uBody.Metadata.Version - 100
	_ = Version
	return nil
}

func (w WatchFlowRun) artifactConnect() error {
	client := sse.NewClient(fmt.Sprintf("%s/watch?resource=flowrun?version=%d", common.EchoerUrl, Version))
	fmt.Printf("%s/watch?resource=flowrun?version=%d", common.EchoerUrl, Version)
	err := client.SubscribeRaw(func(msg *sse.Event) {
		var a FlowRun
		err := json.Unmarshal(msg.Data, &a)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%v\n", a.Metadata)
		//b, _ := strconv.Atoi(a["metadata"]["version"])
		//fmt.Println(b)
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var a WatchFlowRun
	if err := a.firstConnect(); err != nil {
		panic(err)
	}
	if err := a.artifactConnect(); err != nil {
		panic(err)
	}
}
