package myhttp

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type Listener interface {
	Accept() (Conn, error)
	Close() error
}

type NetSocket struct {
	fd       int
	remoteIP net.IP
}

func NewNetSocket(addr string) *NetSocket {

	syscall.ForkLock.Lock()
	listenFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	checkErr(err) // normally, I would not use such a function to wrap error handling,
	// however here we use a lot of log.Fatal and it can save a lot of code.
	syscall.ForkLock.Unlock()

	checkErr(syscall.SetsockoptInt(listenFd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1))

	sp := strings.Split(addr, ":")
	port := 80
	if len(sp) > 1 {
		port, err = strconv.Atoi(sp[1])
		checkErr(err)
	}

	sa := &syscall.SockaddrInet4{Port: port}
	copy(sa.Addr[:], sp[0])

	checkErr(syscall.Bind(listenFd, sa))

	checkErr(syscall.Listen(listenFd, syscall.SOMAXCONN))

	return &NetSocket{listenFd, nil}
}

func (ns *NetSocket) Accept() (Conn, error) {

	nfd, rsa, err := rawAccept(ns.fd)
	if err == nil {
		syscall.CloseOnExec(nfd)
	}

	if err != nil {
		return nil, err
	}

	sa, err := parseIP(rsa)
	if err != nil {
		return nil, err
	}
	ip := net.ParseIP(convertIP(sa.Addr))

	return &NetSocket{nfd, ip}, nil
}

func (ns *NetSocket) Close() error {
	return syscall.Close(ns.fd)
}

func (ns *NetSocket) Write(p []byte) (int, error) {
	n, err := syscall.Write(ns.fd, p)
	if err != nil {
		return 0, nil
	}
	return n, nil
}

func (ns *NetSocket) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	n, err := syscall.Read(ns.fd, p)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (ns *NetSocket) Addr() net.IP {
	return ns.remoteIP
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func accept(s int, rsa *syscall.RawSockaddr, addrlen *uint32) (fd int, err error) {
	r0, _, e1 := syscall.Syscall(syscall.SYS_ACCEPT, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)))
	fd = int(r0)
	if e1 != 0 {
		err = errors.New("accept error")
	}
	return
}

func rawAccept(fd int) (int, *syscall.RawSockaddr, error) {
	var rsa syscall.RawSockaddr
	var len uint32 = 0x70 // syscall.SizeofSockaddrAny
	nfd, err := accept(fd, &rsa, &len)
	return nfd, &rsa, err
}

func convertIP(data [4]byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", int(data[0]), int(data[1]), int(data[2]), int(data[3]))
}

func parseIP(rsa *syscall.RawSockaddr) (*syscall.SockaddrInet4, error) {
	pp := (*syscall.RawSockaddrInet4)(unsafe.Pointer(rsa))
	sa := new(syscall.SockaddrInet4)
	p := (*[2]byte)(unsafe.Pointer(&pp.Port))
	sa.Port = int(p[0])<<8 + int(p[1])
	for i := 0; i < len(sa.Addr); i++ {
		sa.Addr[i] = pp.Addr[i]
	}
	return sa, nil
}
