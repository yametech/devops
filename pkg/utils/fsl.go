package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/go-flowrun"
	"time"
)

type ArtifactCIInfo struct {
	Branch      string `json:"branch"`
	CodeType    string `json:"codeType"`
	CommitID    string `json:"commitId"`
	GitUrl      string `json:"gitUrl"`
	OutPut      string `json:"output"`
	ProjectFile string `json:"projectFile"`
	ProjectPath string `json:"projectPath"`
	RetryCount  uint   `json:"retryCount"`
	ServiceName string `json:"serviceName"`
	appID       string
}

func (a ArtifactCIInfo) SendEchoer() error {
	if a.appID == "" {
		return errors.New("UUID is not none")
	}

	flowRun := &flowrun.FlowRun{
		EchoerUrl: common.EchoerUrl,
		Name:      fmt.Sprintf("%s_%d", common.DefaultNamespace, time.Now().UnixNano()),
	}
	flowRunStep := map[string]string{
		"SUCCESS": "done", "FAIL": "done",
	}
	flowRunAction, err := core.ToMap(a)
	if err != nil {
		return err
	}

	flowRunStepName := fmt.Sprintf("%s_%s", common.EchoerCI, a.appID)
	flowRun.AddStep(flowRunStepName, flowRunStep, common.EchoerCI, flowRunAction)

	flowRunData := flowRun.Generate()
	fmt.Println(flowRunData)
	//if !flowRun.Create(flowRunData) {
	//	return errors.New("send fsm error")
	//}
	return nil
}

//func main() {
//	a := ArtifactCIInfo{
//		Branch: "env-dev",
//		CodeType: "java-maven",
//		CommitID: "realtime-compute-service-845b6f8e99efa6d773d84135395a361f04b5638b",
//		GitUrl: "http://git.ym/g-infrast/ocean-realtime-compute-service.git",
//		OutPut: "registry-d.ym/g-infrast",
//		ProjectFile: "",
//		ProjectPath: "",
//		RetryCount: 15,
//		ServiceName: "realtime-compute-service",
//	}
//	//b, _ := core.ToMap(a)
//	//c, _ := json.Marshal(b)
//	//fmt.Println(b)
//	//fmt.Println(string(c))
//	if err := a.SendEchoer();err != nil {
//		panic(err)
//	}
//}
