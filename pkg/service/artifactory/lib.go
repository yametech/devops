package artifactory

import (
	"encoding/json"
	"fmt"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/resource/artifactory"
	"github.com/yametech/go-flowrun"
	"strconv"
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

func MergeEnvVar(a *artifactory.Container) {
	env := &a.Environment
	m := make(map[string]interface{})
	for _, k := range *env {
		m[Strval(k["name"])] = k["envvalue"]
	}
	*env = make([]map[string]interface{}, 0)
	for k, v := range m {
		*env = append(*env, map[string]interface{}{
			"name":     k,
			"envvalue": v,
		})
	}
}

//interface转string
func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
