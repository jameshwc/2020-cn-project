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
	request, err := parseRequest(buf)
	if err != nil {
		log.Fatal(err)
	}

	msg := []byte("HTTP/1.0 200 OK\r\nContent-Type: text/html\r\n\r\n")
	n, err = conn.Write(msg)
	if err != nil {
		return
	}
	serve(conn, request)
}

func serve(conn net.Conn, req *Request) {
	beforeServe(conn, req)
	switch req.URI {
	case "/":
		homepage(conn, req)
	}
	conn.Write([]byte(header2string(req.Headers)))
	afterServe(conn, req)
}
