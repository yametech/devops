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
	GlobalConfig     = "globalconfig"
)

var (
	EchoerCI = "artifactoryCI"
	EchoerCD = "artifactoryCD"
	GitUser  = "git_user"
	GitPW    = "git_password"
)

func init() {
	flag.StringVar(&EchoerCI, "echoerci", EchoerCI, "-echoerci=artifactoryCI")
	flag.StringVar(&EchoerCD, "echoercd", EchoerCD, "-echoercd=artifactoryCD")
	flag.StringVar(&GitUser, "gituser", GitUser, "-gituser=git_user")
	flag.StringVar(&GitPW, "gitpw", GitPW, "-gitpw=git_password")

}
