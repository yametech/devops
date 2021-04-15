package artifactory

import "github.com/yametech/devops/pkg/core"

type RequestArtifactory struct {
	GitUrl    string `json:"git_url" bson:"git_url"`
	AppName   string `json:"app_name" bson:"app_name"`
	Branch    string `json:"branch" bson:"branch"`
	Tag       string `json:"tag" bson:"tag"`
	Remarks   string `json:"remarks" bson:"remarks"`
	Language  string `json:"language" bson:"language"`
	ImagesHub string `json:"images_hub" bson:"images_hub"`
}

type ArtifactorySpec struct {
	ImagesUrl string `json:"images_url" bson:"images_url"`
	GitUrl    string `json:"git_url" bson:"git_url"`
	AppName   string `json:"app_name" bson:"app_name"`
	Branch    string `json:"branch" bson:"branch"`
	Tag       string `json:"tag" bson:"tag"`
	Remarks   string `json:"remarks" bson:"remarks"`
	Language  string `json:"language" bson:"language"`
	ImagesHub string `json:"images_hub" bson:"images_hub"`
}

type ResArtifactory struct {
	core.Metadata `json:"metadata"`
	Spec          ArtifactorySpec `json:"spec"`
}
