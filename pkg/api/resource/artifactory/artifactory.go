package artifactory

import "github.com/yametech/devops/pkg/core"

type RequestArtifactory struct {
	GitUrl    string `json:"git_url"`
	AppName   string `json:"app_name"`
	Branch    string `json:"branch"`
	Tag       string `json:"tag"`
	Remarks   string `json:"remarks"`
	Language  string `json:"language"`
	ImagesHub string `json:"images_hub"`
}

type ArtifactorySpec struct {
	ImagesUrl string `json:"images_url"`
	GitUrl    string `json:"git_url"`
	AppName   string `json:"app_name"`
	Branch    string `json:"branch"`
	Tag       string `json:"tag"`
	Remarks   string `json:"remarks"`
	Language  string `json:"language"`
	ImagesHub string `json:"images_hub"`
}

type ResArtifactory struct {
	core.Metadata `json:"metadata"`
	Spec          ArtifactorySpec `json:"spec"`
}
