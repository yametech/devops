package appservice_controller

import (
	"flag"
	"github.com/yametech/devops/pkg/controller"
	"github.com/yametech/devops/pkg/store/mongo"
)

var storageUri string

func main() {
	flag.StringVar(&storageUri, "storage_uri", "mongodb://10.200.10.46:27017/devops", "127.0.0.1:3306")
	flag.Parse()

	store, err, errC := mongo.NewMongo(storageUri)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := controller.NewPipelineController(store).Run(); err != nil {
			errC <- err
		}
	}()

	if e := <-errC; e != nil {
		panic(e)
	}
}
