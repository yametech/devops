package resource

import (
	"github.com/yametech/devops/pkg/core"
)

type RequestUser struct {
	Name     string                 `json:"name"`
	Kind     string                 `json:"kind"`
	Password string                 `json:"password"`
	NickName string                 `json:"nick_name"`
	Labels   map[string]interface{} `json:"labels"`
}

type UserSpec struct {
	LastLogin int64  `json:"last_login" bson:"last_login"`
	NickName  string `json:"nick_name" bson:"nick_name"`
}

type RespUser struct {
	core.Metadata `json:"metadata"`
	Spec          UserSpec `json:"spec"`
}
