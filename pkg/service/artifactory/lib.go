package artifactory

import (
	"fmt"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/go-flowrun"
	"time"
	"unicode"
)

func SendEchoer(stepName string, actionName string, a map[string]interface{}) bool {
	flowRun := &flowrun.FlowRun{
		EchoerUrl: common.EchoerUrl,
		Name:      fmt.Sprintf("%s_%d", common.DefaultNamespace, time.Now().UnixNano()),
	}
	flowRunStep := map[string]string{
		"SUCCESS": "done", "FAIL": "done",
	}

	flowRun.AddStep(stepName, flowRunStep, actionName, a)

	flowRunData := flowRun.Generate()
	fmt.Println(flowRunData)
	return flowRun.Create(flowRunData)
}

func IsChinese(str string) bool {
	var count int
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0

}
