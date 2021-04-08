package resource

import "github.com/yametech/devops/pkg/core"

type ArtifactSpec struct {
	Project   string `json:"project" bson:"project" gorm:"size:255;default:''"`
	Version   string `json:"version" bson:"version" gorm:"size:255;default:''"`
	GitURL    string `json:"git_url" bson:"git_url" gorm:"size:255;default:''"`
	Branch    string `json:"branch" bson:"branch" gorm:"size:255;default:''"`
	Tag       string `json:"tag" bson:"tag" gorm:"size:255;default:''"`
	Language  string `json:"language" bson:"language" gorm:"size:255;default:''"`
	ImagesURL string `json:"images_url" bson:"images_url" gorm:"size:255;default:''"`
	AppCode   string `json:"app_code" bson:"app_code" gorm:"size:255;default:''"`
	AppName   string `json:"app_name" bson:"app_name" gorm:"size:255;default:''"`
	Detail    string `json:"detail" bson:"detail" gorm:"size:255;default:''"`
	Package   string `json:"package" bson:"package" gorm:"size:255;default:''"`
	Md5       string `json:"md5" bson:"md5" gorm:"size:255;default:''"`

	CreateUserID string `json:"create_user_id" bson:"create_user_id" gorm:"size:255;not null"`
	ProjectFile  string `json:"projectFile" bson:"project_file" gorm:"size:255;default:''"`
	ProjectPath  string `json:"projectPath" bson:"project_path" gorm:"size:255;default:''"`
}

type Artifact struct {
	Metadata core.Metadata `json:"metadata" gorm:"embedded"`
	Spec     ArtifactSpec  `json:"spec" gorm:"embedded"`
}
