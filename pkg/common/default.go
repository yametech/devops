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
	ChildModule      = "childmodule"
	CollectionModule = "collectionmodule"
	AllModule        = "allmodule"
	ModuleEntry      = "moduleentry"
	RecentVisit      = "recentvisit"
	Topmodule        = "topmodule"
)

var (
	//FlowRun Controller
	EchoerCI = "artifactoryCI"

	//Artifactory
	SmartCityCD      = "artifactoryCD"
	AzureCD          = "artifactoryAzureCD"
	TungChungCD      = "artifactoryTungChungCD"
	GitUser          = "gituser"
	GitPW            = "gitpw"
	RegistryUser     = "username"
	RegistryPW       = "password"
	StoreImagesCount = 5

	//gateway etcd
	EtcdAddress = "http://10.200.65.200:2379"
)

func init() {
	flag.StringVar(&EchoerCI, "echoerci", EchoerCI, "-echoerci=artifactoryCI")
	flag.StringVar(&SmartCityCD, "smartcitycd", SmartCityCD, "-smartcitycd=artifactoryCD")
	flag.StringVar(&AzureCD, "azurecd", AzureCD, "-azurecd=artifactoryAzureCD")
	flag.StringVar(&TungChungCD, "tungchungcd", TungChungCD, "-tungchungcd=artifactoryAzureCD")
	flag.StringVar(&GitUser, "gituser", GitUser, "-gituser=git_user")
	flag.StringVar(&GitPW, "gitpw", GitPW, "-gitpw=git_password")
	flag.StringVar(&RegistryUser, "registryuser", RegistryUser, "-registryuser=registry_user")
	flag.StringVar(&RegistryPW, "registrypw", RegistryPW, "-registrypw=registry_password")
	flag.StringVar(&EtcdAddress, "etcdaddress", EtcdAddress, "-etcdaddress=etcdaddress")
	flag.IntVar(&StoreImagesCount, "storeimgescount", StoreImagesCount, "-storeimagescount=5")
}
