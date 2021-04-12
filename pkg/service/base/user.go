package base

import (
	"fmt"
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

func (u *User) List(page, pageSize int) (*[]resource.User, int64, error) {
	user := make([]resource.User, 0)
	offset := (page - 1) * pageSize
	count, err := u.IService.List("users", offset, pageSize, false, &user)
	if err != nil {
		return nil, count, err
	}
	return &user, count, nil

}

func (u *User) Query(filter map[string]interface{}, page, pageSize int) (*[]resource.User, int64, error) {
	user := make([]resource.User, 0)
	offset := (page - 1) * pageSize
	count, err := u.IService.Query("users", filter, offset, pageSize, false, &user)
	if err != nil {
		return nil, count, err
	}

	return &user, count, nil
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
	count, err := u.IService.Query("users", map[string]interface{}{"uuid": uuid}, -1, -1, false, user)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("user不存在")
	}
	err = u.IService.Update(user, request)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Delete(uuid string) error {
	user := &resource.User{} // "users", filter, offset, pageSize, false, &base
	count, err := u.IService.Query("users", map[string]interface{}{"uuid": uuid}, -1, -1, false, user)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user不存在")
	}
	user.IsDelete = true
	err = u.IService.Save(user)
	if err != nil {
		return err
	}
	return nil

}
