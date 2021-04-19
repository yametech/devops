package base

import (
	apiResource "github.com/yametech/devops/pkg/api/resource"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource"
	"github.com/yametech/devops/pkg/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	service.IService
}

func NewUser(i service.IService) *UserService {
	return &UserService{i}
}

func (u *UserService) Watch(version string) (chan core.IObject, chan struct{}) {
	objectChan := make(chan core.IObject, 32)
	closed := make(chan struct{})
	u.IService.Watch(common.DefaultNamespace, common.User, string(resource.UserKind), version, objectChan, closed)
	return objectChan, closed
}

func (u *UserService) List(name string, page, pageSize int64) ([]interface{}, int64, error) {
	offset := (page - 1) * pageSize
	filter := map[string]interface{}{}
	if name != "" {
		filter["metadata.name"] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + name + ".*", Options: "i"}}
	}
	sort := map[string]interface{}{
		"metadata.version": -1,
	}

	data, err := u.IService.ListByFilter(common.DefaultNamespace, common.User, filter, sort, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.IService.Count(common.DefaultNamespace, common.User, filter)
	if err != nil {
		return nil, 0, err
	}
	return data, count, nil
}

func (u *UserService) Create(reqUser *apiResource.RequestUser) error {
	user := &resource.User{
		Metadata: core.Metadata{
			Name:   reqUser.Name,
			Kind:   reqUser.Kind,
			Labels: reqUser.Labels,
		},
		Spec: resource.UserSpec{
			Password: reqUser.Password,
			NickName: reqUser.NickName,
		},
	}
	user.GenerateVersion()
	_, err := u.IService.Create(common.DefaultNamespace, common.User, user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) GetByUUID(uuid string) (*resource.User, error) {
	user := &resource.User{}
	err := u.IService.GetByUUID(common.DefaultNamespace, common.User, uuid, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) Update(uuid string, reqUser *apiResource.RequestUser) (core.IObject, bool, error) {
	user := &resource.User{
		Metadata: core.Metadata{
			Name:   reqUser.Name,
			Kind:   reqUser.Kind,
			Labels: reqUser.Labels,
		},
		Spec: resource.UserSpec{
			Password: reqUser.Password,
			NickName: reqUser.NickName,
		},
	}
	user.GenerateVersion()
	return u.IService.Apply(common.DefaultNamespace, common.User, uuid, user)
}

func (u *UserService) Delete(uuid string) error {
	err := u.IService.Delete(common.DefaultNamespace, common.User, uuid)
	if err != nil {
		return err
	}
	return nil
}
