package myhttp

import (
	"bytes"
	"errors"
	"net/url"
)

type Request struct {
	Method   string
	URL      *url.URL
	Version  string //http
	Headers  map[string][]string
	Body     string
	PostForm url.Values
}

var (
	ErrRequestFormatError     = errors.New("Request not follow HTTP/1.1 protocol (RFC 2616)")
	ErrRequestMethodIncorrect = errors.New("Request HTTP method incorrect")
	ErrRequestHeaderKeyValue  = errors.New("Request headers not follow {$key: $value} format")
	ErrRequestParseURLError   = errors.New("Request parse url error")
)

func parseRequest(dat []byte) (*Request, error) {

	split := bytes.SplitN(dat, []byte{'\r', '\n', '\r', '\n'}, 2)
	if len(split) < 2 {
		return nil, ErrRequestFormatError
	}

	method, uri, version, err := parseRequestLine(split[0])
	if err != nil {
		return nil, err
	}

	headers, err := parseHeaders(split[0])
	if err != nil {
		return nil, err
	}

	body := string(split[1])

	url, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, ErrRequestParseURLError
	}

	return &Request{method, url, version, headers, body, nil}, nil
}

func parseRequestLine(dat []byte) (method, uri, version string, err error) {
	split := bytes.SplitN(dat, []byte{' '}, 3)
	method = string(split[0])
	switch method {
	case "GET", "POST", "HEAD", "PUT", "DELETE", "OPTIONS", "CONNECT", "TRACE":
	default:
		err = ErrRequestMethodIncorrect
		return
	}
	uri = string(split[1]) // TODO: check?
	version = string(split[2])
	return
}

func parseHeaders(dat []byte) (h map[string][]string, err error) {
	h = make(map[string][]string)
	headers := bytes.Split(dat, []byte{'\r', '\n'})

	for i := range headers {
		if i == 0 { // skip request line
			continue
		}

		s := bytes.SplitN(headers[i], []byte{':', ' '}, 2)
		if len(s) < 2 {
			err = ErrRequestHeaderKeyValue
			return
		}

		key := string(s[0])
		val := string(s[1])

		if _, ok := h[key]; !ok {
			h[key] = make([]string, 0)
		}

		h[key] = append(h[key], val)
	}
	return
}

func (r *Request) ParseForm() (err error) {
	if r.PostForm == nil {
		r.PostForm, err = url.ParseQuery(r.Body)
	}
	return
}
