package main

import (
	"flag"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/api/action/globalconfigservice"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/store/mongo"
)

var storageUri string

func main() {
	flag.StringVar(&storageUri, "storage_uri", "mongodb://127.0.0.1:27017/admin", "-storage_uri=mongodb://127.0.0.1:27017/admin")

	flag.Parse()

	store, err, errC := mongo.NewMongo(storageUri)
	if err != nil {
		panic(err)
	}

	baseService := service.NewBaseService(store)
	server := api.NewServer(baseService)

	globalconfigservice.NewGlobalServiceServer("globalconfig", server)

	go func() {
		if err := server.Run(":8080"); err != nil {
			errC <- err
		}
	}()

	if e := <-errC; e != nil {
		panic(e)
	}
}
