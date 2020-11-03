package main

import (
	"fmt"
	"website/conf"
	"website/log"
	"website/model/mongo"
	"website/myhttp"
	"website/view"
	viewuser "website/view/user"

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

	listen := myhttp.NewNetSocket(endPoint)
	defer listen.Close()

	router := myhttp.NewRouter()

	router.GET("/", view.Homepage)
	router.POST("/messages", view.AddMessage)
	router.GET("/messages", view.ShowMessage)
	router.GET("/register", viewuser.RegisterForm)
	router.POST("/register", viewuser.Register)
	router.GET("/login", viewuser.LoginForm)
	router.POST("/login", viewuser.Login)
	router.GET("/logout", viewuser.Logout)
	router.GET("/video/video.mp4", view.Video)
	router.GET("/video_demo", view.VideoDemo)

	for {
		conn, errs := listen.Accept()
		if errs != nil {
			fmt.Println("accept failed")
			continue
		}
		go myhttp.Handle(conn, router)
	}
}
