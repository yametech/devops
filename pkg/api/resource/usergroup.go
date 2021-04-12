package resource

import "github.com/yametech/devops/pkg/core"

type RequestUserProject struct {
	Name   string                 `json:"name" bson:"name"`
	Kind   string                 `json:"kind"  bson:"kind"`
	Labels map[string]interface{} `json:"labels" bson:"labels"`

	ProjectFile  string `json:"project_file" bson:"project_file" `
	ProjectPath  string `json:"project_path" bson:"project_path" `
	CreateUserID string `json:"create_user_id" bson:"create_user_id" `
}

type UserProjectSpec struct {
	CreateUserID string   `json:"create_user_id" bson:"create_user_id" `
	ProjectFile  string   `json:"projectFile" bson:"project_file" `
	ProjectPath  string   `json:"projectPath" bson:"project_path" `
	CreateUser   RespUser `json:"create_user"`
}

type RespUserProject struct {
	core.Metadata `json:"metadata" `
	Spec          UserProjectSpec `json:"spec" `
}
