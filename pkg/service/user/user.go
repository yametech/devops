package user

import (
	apiResource "github.com/yametech/devops/pkg/api/resource"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource"
	"github.com/yametech/devops/pkg/service"
)

type User struct {
	service.IService
}

func NewUser(i service.IService) *User {
	return &User{i}
}

func (u *User) List() (*[]resource.User, error) {
	user := make([]resource.User, 0)
	err := u.IService.List(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (u *User) Query(filter map[string]interface{}) (*[]resource.User, error) {
	user := make([]resource.User, 0)
	err := u.IService.Query(filter, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) Create(request *apiResource.RequestUser) (*resource.User, error) {
	user := &resource.User{
		Metadata: core.Metadata{
			Name: request.Name,
			Kind: request.Kind,
		},
		Spec: resource.UserSpec{
			NickName: request.NickName,
			Username: request.Username,
			Password: request.Password,
		},
	}
	user.GenerateVersion()
	err := u.IService.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Update(uuid string, request map[string]interface{}) (*resource.User, error) {
	user := &resource.User{}
	err := u.IService.Query(map[string]interface{}{"uuid": uuid}, user)
	if err != nil {
		return nil, err
	}
	err = u.IService.Update(user, request)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Delete(uuid string) error {
	user := &resource.User{}
	err := u.IService.Query(map[string]interface{}{"uuid": uuid}, user)
	if err != nil {
		return err
	}
	user.IsDelete = true
	err = u.IService.Save(user)
	if err != nil {
		return err
	}
	return nil

}
