package main

import (
	"flag"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/api/action/artifactory"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/store/mongo"
	"github.com/yametech/go-insect"
)

var storageUri string

func main() {
	flag.StringVar(&storageUri, "storage_uri", "mongodb://10.200.10.46:27017/admin", "127.0.0.1:3306")
	flag.Parse()

	store, err, errC := mongo.NewMongo(storageUri)
	if err != nil {
		panic(err)
	}

	baseService := service.NewBaseService(store)
	server := api.NewServer(baseService)

	//注册网关
	go_insect.GlobalEtcdAddress = common.EtcdAddress
	go_insect.INSECT_SERVER_PORT = 8080
	go_insect.INSECT_SERVER_NAME = "artifactory"
	go go_insect.EtcdProxy()

	artifactory.NewArBaseServer("artifactory", server)

	go func() {
		if err := server.Run(":8080"); err != nil {
			errC <- err
		}
	}()

	if e := <-errC; e != nil {
		panic(e)
	}
}
