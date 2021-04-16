package resource

import (
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/store"
	"github.com/yametech/devops/pkg/store/gtm"
)

const UserKind core.Kind = "user"

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

// User impl Coder
func (*User) Decode(op *gtm.Op) (core.IObject, error) {
	user := &User{}
	if err := core.ObjectToResource(op.Data, user); err != nil {
		return nil, err
	}
	return user, nil
}

func init() {
	store.AddResourceCoder(string(UserKind), &User{})
}
