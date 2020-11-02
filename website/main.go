package main

import (
	"fmt"
	"website/conf"
	"website/log"
	"website/model/mongo"
	"website/myhttp"
	"website/view"

	_ "github.com/joho/godotenv/autoload"
)

func init() {
	conf.Setup()
	log.Setup()
	mongo.Setup()
}

func main() {
	fmt.Println("start server....")
	endPoint := fmt.Sprintf(":%d", conf.ServerConfig.HttpPort)
	listen := myhttp.NewListener(endPoint)
	router := myhttp.NewRouter()
	router.GET("/", view.Homepage)
	router.POST("/messages", view.AddMessage)
	router.GET("/messages", view.ShowMessage)
	for {
		conn, errs := listen.Accept()
		if errs != nil {
			fmt.Println("accept failed")
			continue
		}
		go myhttp.Handle(conn, router)
	}
}
