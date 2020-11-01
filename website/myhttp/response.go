package myhttp

import (
	"net"
	"strconv"
	"strings"
)

type Responser interface {
	Write([]byte) (int, error)
}

type Response struct {
	StatusCode   int
	Body         string
	Header       map[string][]string
	ReasonPhrase string
	Size         int
	rawData      []byte
	ErrorCode    string // cloudflare
	conn         net.Conn
}

func NewResponse(conn net.Conn) *Response {
	return &Response{conn: conn}
}
func (r *Response) Write(conn net.Conn) {
	statusLine := r.statusLine()
	headers := ""

	for k, v := range r.Header {
		for _, vv := range v {
			headers = headers + k + ": " + vv + "\r\n"
		}
	}

	r.conn.Write([]byte(strings.Join([]string{statusLine, headers, r.Body}, "\r\n")))
}

func (r *Response) statusLine() string {
	return "HTTP/1.0 " + strconv.Itoa(r.StatusCode) + r.ReasonPhrase
}
