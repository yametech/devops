package common

import "flag"

const (
	DefaultNamespace = "devops"
	User             = "user"
	UserProject      = "userproject"
	EchoerUrl        = "http://10.200.65.192:8080"
	Artifactory      = "artifactory"
	AppConfig        = "appconfig"
	AppProject       = "appproject"
	Namespace        = "namespace"
	History          = "history"
	AppResource      = "appresource"
	GlobalConfig     = "globalconfig"
	WorkOrder        = "workorder"
)

var (
	EchoerCI     = "artifactoryCI"
	EchoerCD     = "artifactoryCD"
	GitUser      = "gituser"
	GitPW        = "gitpw"
	RegistryUser = "username"
	RegistryPW   = "password"
)

func init() {
	flag.StringVar(&EchoerCI, "echoerci", EchoerCI, "-echoerci=artifactoryCI")
	flag.StringVar(&EchoerCD, "echoercd", EchoerCD, "-echoercd=artifactoryCD")
	flag.StringVar(&GitUser, "gituser", GitUser, "-gituser=git_user")
	flag.StringVar(&GitPW, "gitpw", GitPW, "-gitpw=git_password")
	flag.StringVar(&RegistryUser, "registryuser", GitUser, "-registryuser=registry_user")
	flag.StringVar(&RegistryPW, "registrypw", RegistryPW, "-registrypw=registry_password")
}
