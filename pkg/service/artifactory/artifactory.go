package artifactory

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/artifactory"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	arResource "github.com/yametech/devops/pkg/resource/artifactory"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	urlPkg "net/url"
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
		"metadata.created_time": -1,
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

	//gitName = https://github.com/yametech/devops.git --> devops
	//gitDirectory = https://github.com/yametech/devops.git --> yametech
	gitName := ""
	gitDirectory := ""
	if strings.Contains(reqAr.GitUrl, "/") {
		sliceTemp := strings.Split(reqAr.GitUrl, "/")
		gitDirectory = sliceTemp[len(sliceTemp)-2]
		if strings.Contains(sliceTemp[len(sliceTemp)-1], ".git") {
			gitName = strings.ToLower(strings.Split(sliceTemp[len(sliceTemp)-1], ".git")[0])
		} else {
			gitName = strings.ToLower(sliceTemp[len(sliceTemp)-1])
		}
		gitName = strings.ReplaceAll(gitName, "_", "-")
	}

	//registry = http://harbor.ym --> harbor.ym
	registry := ""
	if strings.Contains(reqAr.Registry, "http://") {
		sliceTemp := strings.Split(reqAr.Registry, "http://")
		registry = sliceTemp[len(sliceTemp)-1]
	} else if strings.Contains(reqAr.Registry, "https://") {
		sliceTemp := strings.Split(reqAr.Registry, "https://")
		registry = sliceTemp[len(sliceTemp)-1]
	} else {
		registry = reqAr.Registry
	}
	projectPath := fmt.Sprintf("%s/%s", registry, gitDirectory)
	imageUrl := fmt.Sprintf("%s/%s", projectPath, gitName)

	if len(reqAr.Tag) == 0 {
		reqAr.Tag = strings.ToLower(utils.NewSUID().String())
	} else {
		reqAr.Tag = strings.ToLower(reqAr.Tag)
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
	err := a.CheckRegistryProject(ar)
	if err != nil {
		return nil, err
	}
	ar.GenerateVersion()
	_, err = a.IService.Create(common.DefaultNamespace, common.Artifactory, ar)
	if err != nil {
		return nil, err
	}
	go a.SendCI(ar, projectPath)
	return ar, nil
}

func (a *ArtifactService) CheckRegistryProject(ar *arResource.Artifact) error {
	var HarborAddress string
	var catalogue string
	if strings.Contains(ar.Spec.Images, "/") {
		SliceTemp := strings.Split(ar.Spec.Images, "/")
		HarborAddress = SliceTemp[0]
		catalogue = SliceTemp[1]
	}
	url := fmt.Sprintf("https://%s/api/v2.0/projects?project_name=%s", HarborAddress, catalogue)
	data := urlPkg.Values{}
	req, err := http.NewRequest("HEAD", url, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(common.RegistryUser, common.RegistryPW)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 30 * time.Second, Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 {
		err = a.CreateRegistryProject(HarborAddress, catalogue)
		if err != nil {
			return err
		}
	} else if resp.StatusCode == 401 {
		return errors.New("Registry 认证失败")
	}
	return nil
}

func (a *ArtifactService) CreateRegistryProject(HarborAddress, projectName string) error {
	url := fmt.Sprintf("https://%s/api/v2.0/projects", HarborAddress)
	body := map[string]interface{}{
		"project_name": projectName,
		"metadata": map[string]interface{}{
			"public": "true",
		},
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(common.RegistryUser, common.RegistryPW)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 30 * time.Second, Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 201 {
		return nil
	}
	return errors.New("构建镜像仓库目录失败！")
}

func (a *ArtifactService) SendCI(ar *arResource.Artifact, output string) {
	arCIInfo := &arResource.ArtifactCIInfo{
		Branch:      ar.Spec.Branch,
		CodeType:    ar.Spec.Language,
		CommitID:    ar.Spec.Tag,
		GitUrl:      ar.Spec.GitUrl,
		OutPut:      output,
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

	stepName := fmt.Sprintf("%s_%s", common.CI, ar.UUID)
	if !SendEchoer(stepName, common.EchoerCI, sendCIInfo) {
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
	b := new(interface{})
	err := a.IService.GetByUUID(common.DefaultNamespace, common.Artifactory, uuid, b)
	if err != nil {
		return err
	}
	err = a.IService.Delete(common.DefaultNamespace, common.Artifactory, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (a *ArtifactService) GetBranch(org string, name string) ([]string, error) {
	url := fmt.Sprintf("http://git.ym/api/v1/repos/%s/%s/branches", org, name)
	req, err := http.NewRequest("GET", url, strings.NewReader(urlPkg.Values{}.Encode()))
	if err != nil {
		panic(err.Error())
	}
	req.SetBasicAuth(common.GitUser, common.GitPW)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 30 * time.Second, Transport: tr}
	res, err := client.Do(req)
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
