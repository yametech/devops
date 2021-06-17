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
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

func (a *ArtifactService) List(name, status string, page, pageSize int64) ([]interface{}, int64, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	if name != "" {
		filter["spec.app_name"] = name
	}
	if status != "" {
		filter["spec.artifact_status"], _ = strconv.Atoi(status)
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
		reqAr.Tag = fmt.Sprintf(
			"%s-%s", reqAr.AppName,
			a.GetCommitByBranch(reqAr.GitUrl, gitDirectory, gitName, reqAr.Branch),
		)
	} else {
		reqAr.Tag = strings.ToLower(reqAr.Tag)
	}

	appName := fmt.Sprintf("%s-%d", reqAr.AppName, time.Now().Unix())

	ar := &arResource.Artifact{
		Metadata: core.Metadata{
			Name: appName,
		},
		Spec: arResource.ArtifactSpec{
			GitUrl:       reqAr.GitUrl,
			AppName:      reqAr.AppName,
			Branch:       reqAr.Branch,
			Tag:          reqAr.Tag,
			Remarks:      reqAr.Remarks,
			Language:     reqAr.Language,
			Registry:     registry,
			ProjectFile:  reqAr.ProjectFile,
			ProjectPath:  reqAr.ProjectPath,
			Images:       imageUrl,
			CreateUserId: reqAr.UserNameID,
			CreateUser:   reqAr.UserName,
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
	go a.HandleRegistryArtifacts(ar)
	go a.CheckImagesCount(reqAr.AppName)
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

	resp, err := connectYm("HEAD", url, nil, []string{common.RegistryUser, common.RegistryPW})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Println("仓库没有该路径，准备创建")
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

	resp, err := connectYm("POST", url, body, []string{common.RegistryUser, common.RegistryPW})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		return nil
	}
	return errors.New("创建镜像仓库目录失败！")
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

	appName := strings.Replace(ar.Spec.AppName, "-", "", -1)
	appName = strings.Replace(appName, "_", "", -1)
	stepName := fmt.Sprintf("%s_%s_%s", common.CI, ar.UUID, appName)
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

func (a *ArtifactService) ExistName(name string) (bool, error) {
	filter := map[string]interface{}{
		"metadata.name": name,
	}
	ar := &arResource.Artifact{}
	err := a.IService.GetByFilter(common.DefaultNamespace, common.Artifactory, ar, filter)
	if err != nil {
		if err.Error() == "notFound" {
			return false, nil
		}
		return false, err
	}
	return true, nil
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
	res, err := connectYm("GET", url, nil, []string{common.GitUser, common.GitPW})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("git连接出错")
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

func (a *ArtifactService) GetBranchByGithub(org string, name string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches", org, name)
	res, err := connectYm("GET", url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("git连接出错")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.Body == nil {
		return nil, nil
	}
	type branch struct {
		Name string `json:"name"`
	}
	var branches []branch
	err = json.Unmarshal(body, &branches)
	if err != nil {
		return nil, err
	}
	sliceBrancher := make([]string, 0)
	for _, v := range branches {
		sliceBrancher = append(sliceBrancher, v.Name)
	}
	return sliceBrancher, nil
}

func (a *ArtifactService) GetCommitByBranch(gitUrl, org, name, branch string) string {
	var (
		url     string
		account []string
	)

	if strings.Contains(gitUrl, "github.com") {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/%s", org, name, branch)
	} else if strings.Contains(gitUrl, "git.ym") {
		url = fmt.Sprintf("http://git.ym/api/v1/repos/%s/%s/branches/%s", org, name, branch)
		account = []string{common.GitUser, common.GitPW}
	} else {
		return strings.ToLower(utils.NewSUID().String())
	}

	res, err := connectYm("GET", url, nil, account)
	if err != nil {
		fmt.Println(err)
		return strings.ToLower(utils.NewSUID().String())
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return strings.ToLower(utils.NewSUID().String())
	}
	if res.Body == nil {
		return strings.ToLower(utils.NewSUID().String())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return strings.ToLower(utils.NewSUID().String())
	}

	var CommitID struct {
		Commit struct {
			ID  string `json:"id"`
			Sha string `json:"sha"`
		} `json:"commit"`
	}

	err = json.Unmarshal(body, &CommitID)
	if err != nil {
		fmt.Println(err)
		return strings.ToLower(utils.NewSUID().String())
	}

	if CommitID.Commit.Sha != "" {
		return CommitID.Commit.Sha
	} else {
		return CommitID.Commit.ID
	}

}

func connectYm(method, url string, body interface{}, account []string) (*http.Response, error) {
	bodyReader := bytes.NewReader([]byte{})
	if body != nil {
		a, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(a)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	if len(account) == 2 {
		req.SetBasicAuth(account[0], account[1])
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := (&http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}).Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *ArtifactService) HandleRegistryArtifacts(ar *arResource.Artifact) {
	SliceImage := strings.Split(ar.Spec.Images, "/")
	if len(SliceImage) != 3 {
		log.Println("cant read artifact images type ", ar.Spec.Images)
		return
	}
	registry := ar.Spec.Registry
	projectName := SliceImage[1]
	repo := SliceImage[2]
	url := fmt.Sprintf("https://%s/api/v2.0/projects/%s/repositories/%s/artifacts", registry, projectName, repo)
	artifacts, err := a.GetRegistryArtifacts(url)
	if err != nil {
		log.Println("get registry Artifacts error, ", err.Error())
		return
	}
	a.deleteArtifacts(url, artifacts)
}

func (a *ArtifactService) GetRegistryArtifacts(url string) (apiResource.RegistryArtifacts, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(common.RegistryUser, common.RegistryPW)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 30 * time.Second, Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	dataByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	artifacts := apiResource.RegistryArtifacts{}
	err = json.Unmarshal(dataByte, &artifacts)
	if err != nil {
		return nil, err
	}
	return artifacts, nil

}

func (a *ArtifactService) deleteArtifacts(url string, artifacts apiResource.RegistryArtifacts) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 30 * time.Second, Transport: tr}
	if len(artifacts) > 5 {
		for _, artifact := range artifacts[5:] {
			reference := artifact.Digest
			deleteURL := fmt.Sprintf("%s/%s", url, reference)
			deleteREQ, _ := http.NewRequest("DELETE", deleteURL, nil)
			deleteREQ.SetBasicAuth(common.RegistryUser, common.RegistryPW)
			_, err := client.Do(deleteREQ)
			if err != nil {
				log.Println("HandleRegistryArtifacts delete artifact error, ", err.Error())
			}
			//fmt.Println(res)
		}
	}

	for _, artifact := range artifacts {
		if artifact.Tags == nil {
			reference := artifact.Digest
			deleteURL := fmt.Sprintf("%s/%s", url, reference)
			deleteREQ, _ := http.NewRequest("DELETE", deleteURL, nil)
			deleteREQ.SetBasicAuth(common.RegistryUser, common.RegistryPW)
			_, err := client.Do(deleteREQ)
			if err != nil {
				log.Println("HandleRegistryArtifacts delete no tag artifact error, ", err.Error())
			}
		}
	}
}

func (a *ArtifactService) CheckImagesCount(appName string) {
	filter := map[string]interface{}{
		"spec.app_name":        appName,
		"spec.artifact_status": arResource.Built, //先考虑删除构建成功的,后面再看怎么清理构建失败的
	}
	sort := map[string]interface{}{
		"metadata.created_time": 1,
	}
	data, err := a.IService.ListByFilter(common.DefaultNamespace, common.Artifactory, filter, sort, 0, 0)
	if err != nil {
		return
	}
	if len(data) > common.StoreImagesCount {
		deleteList := &[]arResource.Artifact{}
		err = utils.UnstructuredObjectToInstanceObj(data[:len(data)-common.StoreImagesCount], deleteList)
		for _, v := range *deleteList {
			a.IService.Delete(common.DefaultNamespace, common.Artifactory, v.Metadata.UUID)
		}
	}
}


func (a *ArtifactService) SyncAr(reqAr *apiResource.SyncArtifact) (string, error) {
	ok, err := a.ExistName(reqAr.Name)
	if err != nil {
		return "检查制品出错", err
	}
	if ok {
		return "已存在制品", nil
	}
	temp := strings.Split(strings.Replace(reqAr.GitUrl, ".git", "", -1), "/")
	appName := strings.ToLower(temp[len(temp)-1])
	ar := &arResource.Artifact{
		Metadata: core.Metadata{
			Name:        reqAr.Name,
			Version:     reqAr.Version,
			CreatedTime: reqAr.CreateTime,
		},
		Spec: arResource.ArtifactSpec{
			ArtifactStatus: arResource.Built,
			Branch:         reqAr.Branch,
			GitUrl:         reqAr.GitUrl,
			Language:       reqAr.Language,
			Tag:            reqAr.Tag,
			Registry:       strings.Split(reqAr.ImagesUrl, "/")[0],
			Images:         fmt.Sprintf("%s/%s", reqAr.ImagesUrl, appName),
			ProjectFile:    reqAr.ProjectFile,
			ProjectPath:    reqAr.ProjectPath,
			AppName:        reqAr.AppName,
			CreateUserId:   reqAr.UserNameID,
			CreateUser:     reqAr.UserName,
			Remarks:        reqAr.Remarks,
		},
	}
	_, err = a.IService.Create(common.DefaultNamespace, common.Artifactory, ar)
	if err != nil {
		return "创建失败", err
	}
	return "创建成功", nil
}

