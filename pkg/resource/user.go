package resource

import (
	"github.com/yoshiakiley/devops-zpk-server/pkg/core"
)

type UserSpec struct {
	Password    string `json:"password" gorm:"size:255;not null"`
	LastLogin   int64  `json:"last_login" bson:"last_login" gorm:""`
	IsSuperuser bool   `json:"is_superuser" bson:"is_superuser" gorm:"type:bool;default:false"`
	Username    string `json:"username" bson:"username" gorm:"size:255"`
	IsStaff     bool   `json:"is_staff" bson:"is_staff" gorm:"type:bool;default:false"`
	NickName    string `json:"nick_name" bson:"nick_name" gorm:"size:255"`
}

type User struct {
	core.Metadata `json:"metadata" gorm:"embedded"`
	Spec          UserSpec `json:"spec" gorm:"embedded"`
}

type RequestUser struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	UUID        string `json:"uuid"`
	Version     string `json:"version"`
	IsDelete    string `json:"is_delete"`
	CreatedTime string `json:"created_time"`
	Labels      string `json:"labels"`
	Password    string `json:"password"`
	LastLogin   string `json:"last_login"`
	IsSuperuser string `json:"is_superuser"`
	Username    string `json:"username"`
	IsStaff     string `json:"is_staff"`
	NickName    string `json:"nick_name"`
}
