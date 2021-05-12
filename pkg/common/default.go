package common

import "flag"

const (
	DefaultNamespace = "devops"
	User             = "user"
	UserProject      = "userproject"

	//Artifactory
	CI          = "CI"
	CD          = "CD"
	EchoerUrl   = "http://10.200.65.192:8080"
	Deploy      = "deploy"
	Artifactory = "artifactory"

	AppConfig    = "appconfig"
	AppProject   = "appproject"
	Namespace    = "namespace"
	History      = "history"
	AppResource  = "appresource"
	GlobalConfig = "globalconfig"
	WorkOrder    = "workorder"

	GlobalModule     = "globalmodule"
	CollectionModule = "collectionmodule"
	AllModule        = "allmodule"
	ModuleEntry      = "moduleentry"
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
)

//func init() {
//	flag.StringVar(&EchoerCI, "echoerci", EchoerCI, "-echoerci=artifactoryCI")
//	flag.StringVar(&SmartCityCD, "smartcitycd", SmartCityCD, "-smartcitycd=artifactoryCD")
//	flag.StringVar(&AzureCD, "azurecd", AzureCD, "-azurecd=artifactoryAzureCD")
//	flag.StringVar(&TungChungCD, "tungchungcd", TungChungCD, "-tungchungcd=artifactoryAzureCD")
//
//	flag.StringVar(&GitUser, "gituser", GitUser, "-gituser=git_user")
//	flag.StringVar(&GitPW, "gitpw", GitPW, "-gitpw=git_password")
//	flag.StringVar(&RegistryUser, "registryuser", GitUser, "-registryuser=registry_user")
//	flag.StringVar(&RegistryPW, "registrypw", RegistryPW, "-registrypw=registry_password")
//}
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
