package artifactory

import (

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/artifactory"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	arResource "github.com/yametech/devops/pkg/resource/artifactory"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	"github.com/yametech/go-flowrun"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"io"
	"net/http"
	urlpkg "net/url"

	"strconv"

	"strings"
	"time"
	"unicode"
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


func (a *ArtifactService) Create(reqAr *apiResource.RequestArtifact) (*arResource.Artifact, error) {
	if IsChinese(reqAr.Branch) || IsChinese(reqAr.Tag) {
		return nil, errors.New("分支和tag不能为中文")
	}

	gitPath := ""
	if strings.Contains(reqAr.GitUrl, "http://") {
		sliceTemp := strings.Split(reqAr.GitUrl, "http://")
		gitPath = sliceTemp[len(sliceTemp)-1]
	} else if strings.Contains(reqAr.GitUrl, "https://") {
		sliceTemp := strings.Split(gitPath, "https://")
		gitPath = sliceTemp[len(sliceTemp)-1]
	}

	gitName := ""
	gitDirectory := ""
	if strings.Contains(gitPath, "/") {
		if sliceTemp := strings.Split(gitPath, "/"); len(sliceTemp) > 2 {
			gitDirectory = sliceTemp[len(sliceTemp)-2]
			gitName = sliceTemp[len(sliceTemp)-1]
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
	imageUrl := fmt.Sprintf("%s/%s", registry, gitName)

	if len(reqAr.Tag) == 0 {
		reqAr.Tag = strings.ToLower(utils.NewSUID().String())
	}

	appName := fmt.Sprintf("%s-%d", reqAr.AppName, time.Now().UnixNano())

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
			Images:      imageUrl,
		},
	}

	ar.GenerateVersion()
	_, err := a.IService.Create(common.DefaultNamespace, common.Artifactory, ar)
	if err != nil {
		return nil, err
	}
	go a.SendCI(ar)
	return ar, nil
}

func (a *ArtifactService) SendCI(ar *arResource.Artifact) {
	arCIInfo := &arResource.ArtifactCIInfo{
		Branch:      ar.Spec.Branch,
		CodeType:    ar.Spec.Language,
		CommitID:    ar.Spec.Tag,
		GitUrl:      ar.Spec.GitUrl,
		OutPut:      ar.Spec.Registry,
		ProjectPath: ar.Spec.ProjectPath,
		ProjectFile: ar.Spec.ProjectFile,
		RetryCount:  15,
		ServiceName: ar.Spec.AppName,
	}
	sendCIInfo, err := core.ToMap(arCIInfo)
	if err != nil {
		ar.Spec.ArtifactStatus = arResource.InitializeFail
		_, _, err = a.IService.Apply(common.DefaultNamespace, common.Artifactory, ar.UUID, ar, false)
		if err != nil {
			fmt.Printf("sendci initialize save error %s", err)
		}
		return
	}
	if !SendEchoer(ar.Metadata.UUID, common.EchoerCI, sendCIInfo) {
		ar.Spec.ArtifactStatus = arResource.InitializeFail
		_, _, err = a.IService.Apply(common.DefaultNamespace, common.Artifactory, ar.UUID, ar, false)
		if err != nil {
			fmt.Printf("sendci sendEchoer fail save error %s", err)
		}
		return
	}
	ar.Spec.ArtifactStatus = arResource.Building
	_, _, err = a.IService.Apply(common.DefaultNamespace, common.Artifactory, ar.UUID, ar, false)
	if err != nil {
		fmt.Printf("sendci sendEchoer success save error %s", err)
	}
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

	flowRunStepName := fmt.Sprintf("%s_%s", actionName, stepName)
	flowRun.AddStep(flowRunStepName, flowRunStep, actionName, a)

	flowRunData := flowRun.Generate()
	fmt.Println(flowRunData)
	if !flowRun.Create(flowRunData) {
		fmt.Println("send fsm error")
		return false
	}
	return true
}


func (a *ArtifactService) GetBanch(org string, name string) ([]string, error) {
	url := fmt.Sprintf("http://git.ym/api/v1/repos/%s/%s/branches", org, name)
	req, err := http.NewRequest("GET", url, strings.NewReader(urlpkg.Values{}.Encode()))
	if err != nil {
		panic(err.Error())
	}
	req.SetBasicAuth(common.GitUser, common.GitPW)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if res != nil {
		var buffer [512]byte
		result := bytes.NewBuffer(nil)
		for {
			n, err := res.Body.Read(buffer[0:])
			result.Write(buffer[0:n])
			if err != nil && err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
		}

		type GetValue struct {
			Name string `json:"name"`
		}

		var uBody []GetValue
		err := json.Unmarshal(result.Bytes(), &uBody)
		if err != nil {
			return nil, err
		}
		sliceBranch := make([]string, 0)
		for _, value := range uBody {
			sliceBranch = append(sliceBranch, value.Name)
		}
		return sliceBranch, err
	}
	return nil, nil
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
		i, err := strconv.Atoi(sliceName[len(sliceName)-1])
		if err != nil {
		return 0
	}
		if i > number {
		number = i
	}
	}
		return number
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

func (a *ArtifactService) GetBanch(org string, name string) ([]string, error) {
	url := fmt.Sprintf("http://git.ym/api/v1/repos/%s/%s/branches", org, name)
	req, err := http.NewRequest("GET", url, strings.NewReader(urlpkg.Values{}.Encode()))
	if err != nil {
		panic(err.Error())
	}
	req.SetBasicAuth(common.GitUser, common.GitPW)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if res != nil {
		var buffer [512]byte
		result := bytes.NewBuffer(nil)
		for {
			n, err := res.Body.Read(buffer[0:])
			result.Write(buffer[0:n])
			if err != nil && err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
		}

		type GetValue struct {
			Name string `json:"name"`
		}

		var uBody []GetValue
		err := json.Unmarshal(result.Bytes(), &uBody)
		if err != nil {
			return nil, err
		}
		sliceBranch := make([]string, 0)
		for _, value := range uBody {
			sliceBranch = append(sliceBranch, value.Name)
		}
		return sliceBranch, err
	}
	return nil, nil
}
