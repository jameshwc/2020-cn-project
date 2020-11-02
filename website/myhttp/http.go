package myhttp

import (
	"fmt"
	"log"
	"net"
)

func Handle(conn net.Conn, r *Router) {
	defer fmt.Println("connection closed...")
	defer conn.Close() // TODO: move it to method?
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return
	}
	request, err := parseRequest(buf)
	if err != nil {
		log.Fatal(err)
	}
	c := NewContext(conn, request)
	serve(c, r)
}

func serve(c *Context, r *Router) {
	beforeServe(c)
	f, ok := r.GetHandler(c.Request.Method, c.Request.URL.Path)
	fmt.Println(c.Request.Method, c.Request.URL.Path)
	if !ok {
		c.NotFound()
		return
	}
	f(c)
	afterServe(c)
}
