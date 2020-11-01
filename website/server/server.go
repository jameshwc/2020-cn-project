package server

import "net"

type Server struct {
	l       net.Listener
	handler map[string]HandleFunc
}

func NewServer(l net.Listener) *Server {
	return &Server{l, make(map[string]HandleFunc)}
}

func (srv *Server) HandlerFunc(f func())