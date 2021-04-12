package resource

import "github.com/yametech/devops/pkg/core"

type UserProjectSpec struct {
	CreateUserID string `json:"create_user_id" bson:"create_user_id" `
	ProjectFile  string `json:"projectFile" bson:"project_file" `
	ProjectPath  string `json:"projectPath" bson:"project_path" `
}

type UserProject struct {
	core.Metadata `json:"metadata" `
	Spec          UserProjectSpec `json:"spec" `
}

func (u *UserProject) Clone() core.IObject {
	result := &UserProject{}
	core.Clone(u, result)
	return result
}
