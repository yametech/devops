package base

import (
	apiResource "github.com/yametech/devops/pkg/api/resource"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
)

type UserProjectService struct {
	service.IService
}

func NewUserProjectService(i service.IService) *UserProjectService {
	return &UserProjectService{i}
}

func (u *UserProjectService) List(page, pageSize int64) ([]*apiResource.RespUserProject, int64, error) {
	offset := (page - 1) * pageSize
	sort := map[string]interface{}{
		"metadata.version": -1,
	}
	unstructured, err := u.IService.List(common.DefaultNamespace, common.UserProject, "", sort, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	userProject, err := u.Structer(unstructured)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.IService.Count(common.DefaultNamespace, common.UserProject, map[string]interface{}{})
	if err != nil {
		return nil, 0, err
	}
	return userProject, count, nil
}

func (u *UserProjectService) Create(reqUserProject *apiResource.RequestUserProject) error {
	userProject := &resource.UserProject{
		Metadata: core.Metadata{
			Name:   reqUserProject.Name,
			Kind:   reqUserProject.Kind,
			Labels: reqUserProject.Labels,
		},
		Spec: resource.UserProjectSpec{
			ProjectFile:  reqUserProject.ProjectFile,
			ProjectPath:  reqUserProject.ProjectPath,
			CreateUserID: reqUserProject.CreateUserID,
		},
	}
	userProject.GenerateVersion()
	_, err := u.IService.Create(common.DefaultNamespace, common.UserProject, userProject)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserProjectService) Structer(unstructured []interface{}) ([]*apiResource.RespUserProject, error) {
	respUserProject := make([]*apiResource.RespUserProject, 0)
	for _, project := range unstructured {
		userProject := apiResource.RespUserProject{}
		err := utils.UnstructuredObjectToInstanceObj(project, &userProject)

		if err != nil {
			return nil, err
		}
		user := apiResource.RespUser{}
		filter := map[string]interface{}{
			"metadata.uuid":      userProject.Spec.CreateUserID,
			"metadata.is_delete": false,
		}
		err = u.GetByFilter(common.DefaultNamespace, common.User, &user, filter)
		if err != nil {
			respUserProject = append(respUserProject, &userProject)
			continue
		}
		userProject.Spec.CreateUser = user
		respUserProject = append(respUserProject, &userProject)
	}
	return respUserProject, nil
}
