package myhttp

import (
	"fmt"
	"log"
	"net"
)

func Handle(conn net.Conn) {
	defer fmt.Println("connection closed...")
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return
	}
	_, err = parseRequest(buf)
	if err != nil {
		log.Fatal(err)
	}

	msg := []byte("HTTP/1.0 200 OK\r\nContent-Type: text/html\r\n\r\n")
	n, err = conn.Write(msg)
	if err != nil {
		return
	}

}
