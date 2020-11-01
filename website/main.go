package main

import (
	"fmt"
	"website/log"
	"website/myhttp"
)

func init() {
	log.Setup()
}

func main() {
	fmt.Println("start server....")
	listen := myhttp.NewListener("0.0.0.0:30006")
	for {
		conn, errs := listen.Accept()
		if errs != nil {
			fmt.Println("accept failed")
			continue
		}
		go myhttp.Handle(conn)
	}
}
