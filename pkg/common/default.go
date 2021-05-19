package common

import "flag"

const (
	DefaultNamespace = "devops"

	//Artifactory
	CI          = "CI"
	CD          = "CD"
	EchoerUrl   = "http://10.200.65.192:8080"
	Deploy      = "deploy"
	Artifactory = "artifactory"

	// Appservice
	AppConfig    = "appconfig"
	AppProject   = "appproject"
	Namespace    = "namespace"
	History      = "appconfigresourcehistory"
	AppResource  = "appresource"
	GlobalConfig = "globalconfig"

	// WorkOrder
	WorkOrder = "workorder"

	// Base
	GlobalModule     = "globalmodule"
	CollectionModule = "collectionmodule"
	AllModule        = "allmodule"
)

var (
	//FlowRun Controller
	EchoerCI = "artifactoryCI"

	//Artifactory
	SmartCityCD  = "artifactoryCD"
	AzureCD      = "artifactoryAzureCD"
	TungChungCD  = "artifactoryTungChungCD"
	GitUser      = "gituser"
	GitPW        = "gitpw"
	RegistryUser = "username"
	RegistryPW   = "password"

	//gateway etcd
	EtcdAddress = "http://10.200.65.200:2379"
)

func init() {
	flag.StringVar(&EchoerCI, "echoerci", EchoerCI, "-echoerci=artifactoryCI")
	flag.StringVar(&SmartCityCD, "smartcitycd", SmartCityCD, "-smartcitycd=artifactoryCD")
	flag.StringVar(&AzureCD, "azurecd", AzureCD, "-azurecd=artifactoryAzureCD")
	flag.StringVar(&TungChungCD, "tungchungcd", TungChungCD, "-tungchungcd=artifactoryAzureCD")
	flag.StringVar(&GitUser, "gituser", "devops-git", "-gituser=git_user")
	flag.StringVar(&GitPW, "gitpw", "devops-git_ABC312", "-gitpw=git_password")
	flag.StringVar(&RegistryUser, "registryuser", "sushaolin", "-registryuser=registry_user")
	flag.StringVar(&RegistryPW, "registrypw", "Ssl19960511.", "-registrypw=registry_password")
}
