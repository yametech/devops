package main

import (
	"flag"
	"github.com/gin-contrib/cors"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/api/action/artifactory"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/store/mongo"
	"time"
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
	server.Engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                            // 允许的前端地址
		AllowMethods:     []string{"PUT", "GET", "DELETE", "POST"}, //允许的方法
		AllowHeaders:     []string{"*"},                            //添加的header
		ExposeHeaders:    []string{"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*" // 允许的前端地址
		},
		MaxAge: 12 * time.Hour,
	}))

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
