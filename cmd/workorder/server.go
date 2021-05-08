package main

import (
	"flag"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/api/action/workorder"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/store/mongo"
)

var storageUri, user, pw, database string

func main() {
	flag.StringVar(&storageUri, "storage_uri", "mongodb://10.200.10.46:27017/devops", "127.0.0.1:3306")
	flag.StringVar(&user, "base", "root", "-base root")
	flag.StringVar(&pw, "pw", "123456", "-pw 123456")
	flag.StringVar(&database, "database", "ccmose", "-database ccmose")
	flag.Parse()

	//errC := make(chan error)
	//store, err := mysql.Setup(storageUri, user, pw, database, errC)
	store, err, errC := mongo.NewMongo(storageUri)

	if err != nil {
		panic(err)
	}

	baseService := service.NewBaseService(store)
	server := api.NewServer(baseService)

	workorder.NewWorkOrder("workorder", server)

	go func() {
		if err := server.Run(":8080"); err != nil {
			errC <- err
		}
	}()

	if e := <-errC; e != nil {
		panic(e)
	}

}
