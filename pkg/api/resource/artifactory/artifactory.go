package artifactory

import (
	"github.com/yametech/devops/pkg/core"
	"time"
)

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
	UserName    string `json:"user_name"`
	UserNameID  string `json:"user_name_id"`
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

type RegistryArtifacts []struct {
	Digest     string      `json:"digest"`
	PullTime   time.Time   `json:"pull_time"`
	PushTime   time.Time   `json:"push_time"`
	References interface{} `json:"references"`
	Tags       []struct {
		ArtifactID   int       `json:"artifact_id"`
		ID           int       `json:"id"`
		Immutable    bool      `json:"immutable"`
		Name         string    `json:"name"`
		PullTime     time.Time `json:"pull_time"`
		PushTime     time.Time `json:"push_time"`
		RepositoryID int       `json:"repository_id"`
		Signed       bool      `json:"signed"`
	} `json:"tags"`
}

type SyncArtifact struct {
	GitUrl      string `json:"git_url"`
	AppName     string `json:"app_name"`
	Branch      string `json:"branch"`
	Tag         string `json:"tag"`
	Remarks     string `json:"remarks"`
	Language    string `json:"language"`
	ImagesUrl   string `json:"images_url"`
	ProjectFile string `json:"project_file"`
	ProjectPath string `json:"project_path"`
	UserName    string `json:"user_name"`
	UserNameID  string `json:"user_name_id"`
	CreateTime  int64  `json:"create_time"`
	Version     int64  `json:"version"`
	Name        string `json:"name"`
}
