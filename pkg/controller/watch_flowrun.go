package main

import (
	"encoding/json"
	"fmt"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/store"
	"io/ioutil"
	"net/http"
)

type WatchFlowRun struct {
	store.IKVStore
}

func (w WatchFlowRun) Run() error {
	panic("implement me")
}

func (w WatchFlowRun) Stop() error {
	panic("implement me")
}

type flowRun struct {
	Spec struct {
		Steps []struct {
			Spec struct {
				RetryCount int `json:"retry_count"`
				Response   struct {
					State string `json:"state"`
				} `json:"response"`
				GlobalVariables interface{} `json:"global_variables"`
				FlowRunUUID     string      `json:"flow_run_uuid"`
				FlowID          string      `json:"flow_id"`
				Data            string      `json:"data"`
				ActionRun       struct {
					ReturnStateMap struct {
						SUCCESS string `json:"SUCCESS"`
						FAIL    string `json:"FAIL"`
					} `json:"return_state_map"`
					Done         bool `json:"done"`
					ActionParams struct {
						ServiceName string `json:"serviceName"`
						RetryCount  int    `json:"retryCount"`
						ProjectPath string `json:"projectPath"`
						ProjectFile string `json:"projectFile"`
						Output      string `json:"output"`
						GitURL      string `json:"gitUrl"`
						CommitID    string `json:"commitId"`
						CodeType    string `json:"codeType"`
						Branch      string `json:"branch"`
					} `json:"action_params"`
					ActionName string `json:"action_name"`
				} `json:"action_run"`
			} `json:"spec"`
			Metadata struct {
				Version int         `json:"version"`
				UUID    string      `json:"uuid"`
				Name    string      `json:"name"`
				Labels  interface{} `json:"labels"`
				Kind    string      `json:"kind"`
			} `json:"metadata"`
		} `json:"steps"`
		LastState      string      `json:"last_state"`
		LastEvent      string      `json:"last_event"`
		LastErr        string      `json:"last_err"`
		HistoryStates  []string    `json:"history_states"`
		GlobalVariable interface{} `json:"global_variable"`
		CurrentState   string      `json:"current_state"`
	} `json:"spec"`
	Metadata struct {
		Version int         `json:"version"`
		UUID    string      `json:"uuid"`
		Name    string      `json:"name"`
		Labels  interface{} `json:"labels"`
		Kind    string      `json:"kind"`
	} `json:"metadata"`
}

//var _ Controller = &WatchFlowRun{}

func (w WatchFlowRun) firstConnect() error {
	resp, err := http.Get(fmt.Sprintf("%s/flowrun/flow_art_1618378932003", common.EchoerUrl))
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var ubody flowRun
	json.Unmarshal(body, &ubody)
	fmt.Println(ubody)
	return nil
}

func main() {
	var a WatchFlowRun
	a.firstConnect()
}
