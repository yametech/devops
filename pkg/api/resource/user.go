package resource

import (
	"github.com/yametech/devops/pkg/core"
)

type RequestUser struct {
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Password string `json:"password"`
	Username string `json:"username"`
	NickName string `json:"nick_name"`
}

type UserSpec struct {
	LastLogin int64  `json:"last_login" bson:"last_login"`
	Username  string `json:"username" bson:"username" gorm:"size:255"`
	NickName  string `json:"nick_name" bson:"nick_name" gorm:"size:255"`
}

type User struct {
	core.Metadata `json:"metadata" gorm:"embedded"`
	Spec          UserSpec `json:"spec" gorm:"embedded"`
}
