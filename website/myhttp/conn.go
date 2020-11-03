package myhttp

import "net"

type Conn interface {
	Read(b []byte) (n int, err error)

	Write(b []byte) (n int, err error)

	Close() error

	Addr() net.IP
}
