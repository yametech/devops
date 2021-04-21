package artifactory

import "github.com/yametech/devops/pkg/core"

type ArtifactStatus uint8

type RequestArtifact struct {
	GitUrl    string `json:"git_url"`
	AppName   string `json:"app_name"`
	Branch    string `json:"branch"`
	Tag       string `json:"tag"`
	Remarks   string `json:"remarks"`
	Language  string `json:"language"`
	ImagesHub string `json:"images_hub"`
}

type ArtifactSpec struct {
	ImagesUrl      string `json:"images_url"`
	GitUrl         string `json:"git_url"`
	AppName        string `json:"app_name"`
	Branch         string `json:"branch"`
	Tag            string `json:"tag"`
	Remarks        string `json:"remarks"`
	Language       string `json:"language"`
	ImagesHub      string `json:"images_hub"`
	ArtifactStatus `json:"artifact_status" bson:"artifact_status"`
}

type RespArtifact struct {
	core.Metadata `json:"metadata"`
	Spec          ArtifactSpec `json:"spec"`
}

type Metadata struct {
	Name        string                 `json:"name" bson:"name"`
	Kind        string                 `json:"kind"  bson:"kind"`
	UUID        string                 `json:"uuid" bson:"uuid"`
	Version     int64                  `json:"version" bson:"version"`
	IsDelete    bool                   `json:"is_delete" bson:"is_delete"`
	CreatedTime int64                  `json:"created_time" bson:"created_time"`
	Labels      map[string]interface{} `json:"labels" bson:"labels"`
}
