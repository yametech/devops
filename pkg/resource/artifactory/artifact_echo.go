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

type ArtifactCDInfo struct {
	ArtifactInfo    map[string]interface{} `json:"artifactInfo"`
	CPULimit        string                 `json:"cpuLimit"`
	CPURequests     string                 `json:"cpuRequests"`
	DeployNamespace string                 `json:"deployNamespace"`
	DeployType      string                 `json:"deployType"`
	MemLimit        string                 `json:"memLimit"`
	MemRequests     string                 `json:"memRequests"`
	Policy          string                 `json:"policy"`
	Replicas        int                    `json:"replicas"`
	ServiceImage    string                 `json:"serviceImage"`
	ServiceName     string                 `json:"serviceName"`
	StorageCapacity string                 `json:"storageCapacity"`
}
