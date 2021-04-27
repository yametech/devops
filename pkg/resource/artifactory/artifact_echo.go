package artifactory

type ArtifactCIInfo struct {
	Branch      string `json:"branch"`
	CodeType    string `json:"codeType"`
	CommitID    string `json:"commitId"`
	GitUrl      string `json:"gitUrl"`
	OutPut      string `json:"output"`
	ProjectFile string `json:"projectFile"`
	ProjectPath string `json:"projectPath"`
	RetryCount  uint   `json:"retryCount"`
	ServiceName string `json:"serviceName"`
}
