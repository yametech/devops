package main

import (
	"flag"
	"github.com/yoshiakiley/devops-zpk-server/pkg/api"
	"github.com/yoshiakiley/devops-zpk-server/pkg/api/action/base"
	"github.com/yoshiakiley/devops-zpk-server/pkg/service"
	"github.com/yoshiakiley/devops-zpk-server/pkg/store/mysql"
)

var storageUri, user, pw, database string

func main() {
	flag.StringVar(&storageUri, "storage_uri", "127.0.0.1:3306", "127.0.0.1:3306")
	flag.StringVar(&user, "user", "root", "-user root")
	flag.StringVar(&pw, "pw", "123456", "-pw 123456")
	flag.StringVar(&database, "database", "ccmose", "-database ccmose")
	flag.Parse()

	errC := make(chan error)
	store, err := mysql.Setup(storageUri, user, pw, database, errC)
	if err != nil {
		panic(err)
	}

	baseService := service.NewBaseService(store)
	server := api.NewServer(baseService)

	base.NewBaseServer("baseserver", server)

	go func() {
		if err := server.Run(":8080"); err != nil {
			errC <- err
		}
	}()

	if e := <-errC; e != nil {
		panic(e)
	}

}
