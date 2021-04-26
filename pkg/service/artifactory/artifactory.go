package artifactory

import (
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	apiResource "github.com/yametech/devops/pkg/api/resource/artifactory"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	arResource "github.com/yametech/devops/pkg/resource/artifactory"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	"github.com/yametech/go-flowrun"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
	"time"
)

type ArtifactService struct {
	service.IService
}

func NewArtifact(i service.IService) *ArtifactService {
	return &ArtifactService{i}
}

func (a *ArtifactService) Watch(version string) (chan core.IObject, chan struct{}) {
	objectChan := make(chan core.IObject, 32)
	closed := make(chan struct{})
	a.IService.Watch(common.DefaultNamespace, common.Artifactory, string(arResource.ArtifactKind), version, objectChan, closed)
	return objectChan, closed
}

func (a *ArtifactService) List(name string, page, pageSize int64) ([]interface{}, int64, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	if name != "" {
		filter["spec.app_name"] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + name + ".*", Options: "i"}}
	}
	sort := map[string]interface{}{
		"metadata.version": -1,
	}

	data, err := a.IService.ListByFilter(common.DefaultNamespace, common.Artifactory, filter, sort, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := a.IService.Count(common.DefaultNamespace, common.Artifactory, filter)
	if err != nil {
		return nil, 0, err
	}
	return data, count, nil

}

func (a *ArtifactService) Create(reqAr *apiResource.RequestArtifact) error {
	gitPath := ""
	if strings.Contains(reqAr.GitUrl, "http://") {
		sliceTemp := strings.Split(reqAr.GitUrl, "http://")
		gitPath = sliceTemp[len(sliceTemp)-1]
	} else if strings.Contains(reqAr.GitUrl, "https://") {
		sliceTemp := strings.Split(gitPath, "https://")
		gitPath = sliceTemp[len(sliceTemp)-1]
	}

	gitDirectory := ""
	if strings.Contains(gitPath, "/") {
		if sliceTemp := strings.Split(gitPath, "/"); len(sliceTemp) > 2 {
			gitDirectory = sliceTemp[len(sliceTemp)-2]
		}
	}

	registry := ""
	if strings.Contains(reqAr.Registry, "http://") {
		sliceTemp := strings.Split(reqAr.Registry, "http://")
		registry = sliceTemp[len(sliceTemp)-1]
	} else if strings.Contains(reqAr.Registry, "https://") {
		sliceTemp := strings.Split(gitPath, "https://")
		registry = sliceTemp[len(sliceTemp)-1]
	}
	registry = fmt.Sprintf("%s/%s", registry, gitDirectory)

	if len(reqAr.Tag) == 0 {
		reqAr.Tag = utils.NewSUID().String()
	}

	appName := fmt.Sprintf("%s-%d", reqAr.AppName, a.GetAppNumber(reqAr.AppName)+1)

	ar := &arResource.Artifact{
		Metadata: core.Metadata{
			Name: appName,
		},
		Spec: arResource.ArtifactSpec{
			GitUrl:      reqAr.GitUrl,
			AppName:     reqAr.AppName,
			Branch:      reqAr.Branch,
			Tag:         reqAr.Tag,
			Remarks:     reqAr.Remarks,
			Language:    reqAr.Language,
			Registry:    registry,
			ProjectFile: reqAr.ProjectFile,
			ProjectPath: reqAr.ProjectPath,
		},
	}

	ar.GenerateVersion()
	_, err := a.IService.Create(common.DefaultNamespace, common.Artifactory, ar)
	if err != nil {
		return err
	}
	//TODO:sendCIEcho
	arCIInfo := &arResource.ArtifactCIInfo{
		Branch:      reqAr.Branch,
		CodeType:    reqAr.Language,
		CommitID:    reqAr.Tag,
		GitUrl:      reqAr.GitUrl,
		OutPut:      registry,
		ProjectPath: reqAr.ProjectPath,
		ProjectFile: reqAr.ProjectFile,
		RetryCount:  15,
		ServiceName: reqAr.AppName,
	}
	sendCIInfo, err := core.ToMap(arCIInfo)
	if err != nil {
		return err
	}
	if SendEchoer(ar.Metadata.UUID, common.EchoerCI, sendCIInfo) {
		//TODO:change CIStatus
		return nil
	}
	return nil
}

func (a *ArtifactService) GetByUUID(uuid string) (*arResource.Artifact, error) {
	ar := &arResource.Artifact{}
	err := a.IService.GetByUUID(common.DefaultNamespace, common.Artifactory, uuid, ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (a *ArtifactService) Update(uuid string, reqAr *apiResource.RequestArtifact) (core.IObject, bool, error) {
	ar := &arResource.Artifact{
		Spec: arResource.ArtifactSpec{
			GitUrl:   reqAr.GitUrl,
			AppName:  reqAr.AppName,
			Branch:   reqAr.Branch,
			Tag:      reqAr.Tag,
			Remarks:  reqAr.Remarks,
			Language: reqAr.Language,
			Registry: reqAr.Registry,
		},
	}
	ar.GenerateVersion()
	return a.IService.Apply(common.DefaultNamespace, common.Artifactory, uuid, ar, false)
}

func (a *ArtifactService) Delete(uuid string) error {
	err := a.IService.Delete(common.DefaultNamespace, common.Artifactory, uuid)
	if err != nil {
		return err
	}
	return nil
}

func SendEchoer(stepName string, actionName string, a map[string]interface{}) bool {
	if stepName == "" {
		fmt.Println("UUID is not none")
		return false
	}

	flowRun := &flowrun.FlowRun{
		EchoerUrl: common.EchoerUrl,
		Name:      fmt.Sprintf("%s_%d", common.DefaultNamespace, time.Now().UnixNano()),
	}
	flowRunStep := map[string]string{
		"SUCCESS": "done", "FAIL": "done",
	}

	flowRunStepName := fmt.Sprintf("%s_%s", actionName, "4444")
	flowRun.AddStep(flowRunStepName, flowRunStep, actionName, a)

	flowRunData := flowRun.Generate()
	fmt.Println(flowRunData)
	if !flowRun.Create(flowRunData) {
		fmt.Println("send fsm error")
		return false
	}
	return true
}

func (a *ArtifactService) GetBranch(gitpath string) ([]string, error) {
	url := fmt.Sprintf("http://%s:%s@%s", common.GitUser, common.GitPW, gitpath)
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:          url,
		SingleBranch: false,
		NoCheckout:   true,
		Depth:        1,
	})
	if err != nil {
		return nil, err
	}

	sliceBranch := make([]string, 0)
	referenceIter, _ := r.References()
	err = referenceIter.ForEach(func(c *plumbing.Reference) error {
		if strings.Contains(string(c.Name()), "refs/remotes/origin/") {
			sliceTemp := strings.Split(string(c.Name()), "refs/remotes/origin/")
			sliceBranch = append(sliceBranch, sliceTemp[len(sliceTemp)-1])
		}
		return nil
	})
	return sliceBranch, err
}

func (a *ArtifactService) GetAppNumber(appName string) int {
	data, _, err := a.List(appName, 1, 0)
	if err != nil {
		return 0
	}
	b, err := json.Marshal(data)
	c := make([]*arResource.Artifact, 0)
	err = json.Unmarshal(b, &c)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	var number = 0
	for _, v := range c {
		sliceName := strings.Split(v.Metadata.Name, "-")
		i, _ := strconv.Atoi(sliceName[len(sliceName)-1])
		if err != nil {
			return 0
		}
		if i > number {
			number = i
		}
	}
	return number
}
