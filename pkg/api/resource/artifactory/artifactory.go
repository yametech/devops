package artifactory

import "github.com/yametech/devops/pkg/core"

type ArtifactStatus uint8

type RequestArtifact struct {
	GitUrl      string `json:"git_url"`
	AppName     string `json:"app_name"`
	Branch      string `json:"branch"`
	Tag         string `json:"tag"`
	Remarks     string `json:"remarks"`
	Language    string `json:"language"`
	Registry    string `json:"registry"`
	ProjectFile string `json:"project_file"`
	ProjectPath string `json:"project_path"`
}

type ArtifactSpec struct {
	ImagesUrl      string `json:"images_url"`
	GitUrl         string `json:"git_url"`
	AppName        string `json:"app_name"`
	Branch         string `json:"branch"`
	Tag            string `json:"tag"`
	Remarks        string `json:"remarks"`
	Language       string `json:"language"`
	Registry       string `json:"registry"`
	ProjectFile    string `json:"project_file"`
	ProjectPath    string `json:"project_path"`
	ArtifactStatus `json:"artifact_status" bson:"artifact_status"`
}

type RespArtifact struct {
	core.Metadata `json:"metadata"`
	Spec          ArtifactSpec `json:"spec"`
}
