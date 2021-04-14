package main

import (
	"flag"
	"github.com/yametech/devops/pkg/store/mongo"
)

var storageUri string

func main() {
	flag.StringVar(&storageUri, "storage_uri", "mongodb://106.13.211.154:27017/admin", "127.0.0.1:3306")
	flag.Parse()

	store, err, errC := mongo.NewMongo(storageUri)

}
