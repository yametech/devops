package resource

import (
	"github.com/yametech/devops/pkg/core"
)

//TODO:and this resource
type UserSpec struct {
	Password    string `json:"password" bson:"password"`
	LastLogin   int64  `json:"last_login" bson:"last_login"`
	IsSuperuser bool   `json:"is_superuser" bson:"is_superuser" `
	IsStaff     bool   `json:"is_staff" bson:"is_staff"`
	NickName    string `json:"nick_name" bson:"nick_name"`
}

type User struct {
	core.Metadata `json:"metadata" `
	Spec          UserSpec `json:"spec"`
}

func (u *User) Clone() core.IObject {
	result := &User{}
	core.Clone(u, result)
	return result
}
