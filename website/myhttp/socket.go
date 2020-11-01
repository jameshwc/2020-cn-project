package myhttp

import (
	"log"
	"net"
)

func NewListener(addr string) net.Listener {
	socket, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	return socket
}
