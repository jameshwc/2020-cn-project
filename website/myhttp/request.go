package myhttp

import (
	"bytes"
	"errors"
	"fmt"
)

type Request struct {
	Headers map[string][]string
}

var (
	ErrHTTPRequestFormatError = errors.New("Request not follow HTTP/1.1 protocol (RFC 2616)")
)

func parseRequest(dat []byte) (*Request, error) {

	split := bytes.SplitN(dat, []byte{'\r', '\n', '\r', '\n'}, 2)
	if len(split) < 2 {
		return nil, ErrHTTPRequestFormatError
	}

	// method, uri, version := parseRequestLine(split[0])

	for _, header := range split {
		fmt.Println(string(header))
	}
	return nil, nil
}

func parseRequestLine(dat []byte) (method, uri, version string) {

}
