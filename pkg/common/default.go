package common

import "flag"

const (
	DefaultNamespace = "devops"
	User             = "user"
	UserProject      = "userproject"
	EchoerUrl        = "http://10.200.65.192:8080"
	AppProject       = "appProject"
	Artifactory      = "artifactory"
	AppConfig        = "appconfig"
)

var (
	EchoerCI = "artifactoryCI"
	EchoerCD = "artifactoryCD"
)

func init() {
	flag.StringVar(&EchoerCI, "echoerci", "artifactoryCI", "-echoerci = artifactoryCI")
	flag.StringVar(&EchoerCD, "echoercd", "artifactoryCD", "-echoercd = artifactoryCD")
}
