package common

import "flag"

const (
	DefaultNamespace = "devops"
	User             = "user"
	UserProject      = "userproject"
	EchoerUrl        = "http://10.200.65.192:8080"
	WorkOrderStatus  = "http://127.0.0.1:8081/workorder/status"
	Artifactory      = "artifactory"
	AppConfig        = "appconfig"
	AppProject       = "appservice"
	ResourcePool     = "resourcepool"
	History          = "history"
	AppResource      = "appresource"
	GlobalConfig     = "globalconfig"
	WorkOrder        = "workorder"
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
