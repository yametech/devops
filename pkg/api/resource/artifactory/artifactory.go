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
	AdditionLinks struct {
		BuildHistory struct {
			Absolute bool   `json:"absolute"`
			Href     string `json:"href"`
		} `json:"build_history"`
	} `json:"addition_links"`
	Digest     string `json:"digest"`
	ExtraAttrs struct {
		Architecture string      `json:"architecture"`
		Author       interface{} `json:"author"`
		Created      time.Time   `json:"created"`
		Os           string      `json:"os"`
	} `json:"extra_attrs"`
	Icon              string      `json:"icon"`
	ID                int         `json:"id"`
	Labels            interface{} `json:"labels"`
	ManifestMediaType string      `json:"manifest_media_type"`
	MediaType         string      `json:"media_type"`
	ProjectID         int         `json:"project_id"`
	PullTime          time.Time   `json:"pull_time"`
	PushTime          time.Time   `json:"push_time"`
	References        interface{} `json:"references"`
	RepositoryID      int         `json:"repository_id"`
	Size              int         `json:"size"`
	Tags              []struct {
		ArtifactID   int       `json:"artifact_id"`
		ID           int       `json:"id"`
		Immutable    bool      `json:"immutable"`
		Name         string    `json:"name"`
		PullTime     time.Time `json:"pull_time"`
		PushTime     time.Time `json:"push_time"`
		RepositoryID int       `json:"repository_id"`
		Signed       bool      `json:"signed"`
	} `json:"tags"`
	Type string `json:"type"`
}


