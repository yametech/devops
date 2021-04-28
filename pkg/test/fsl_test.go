package test

import (
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/artifactory"
	artifactory2 "github.com/yametech/devops/pkg/service/artifactory"
	"testing"
)

func TestCreatCIfsl(t *testing.T) {
	a := &artifactory.ArtifactCIInfo{
		Branch:      "master",
		CodeType:    "web",
		CommitID:    "devopsui-845b6f8e99efa6d773d84135395a361f04b5638b",
		GitUrl:      "http://git.ym/fengdr/Devops-prototype.git",
		OutPut:      "registry-d.ym/fengdr",
		ProjectFile: "",
		ProjectPath: "",
		RetryCount:  15,
		ServiceName: "devopsui",
	}
	b, _ := core.ToMap(a)
	//c, _ := json.Marshal(b)
	//fmt.Println(b)
	//fmt.Println(string(c))
	artifactory2.SendEchoer("", common.EchoerCI, b)
}
