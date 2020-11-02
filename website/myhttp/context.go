package myhttp

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
	"website/log"
)

type Context struct {
	conn          net.Conn
	Headers       Header
	Request       *Request
	isHeaderWrite bool
	queryCache    url.Values
	formCache     url.Values
}

const (
	htmlContentType = "text/html; charset=utf-8"
	jsonContentType = "application/json; charset=utf-8"
)

func NewContext(c net.Conn, req *Request) *Context {
	return &Context{c, make(map[string][]string), req, false, nil, nil}
}

func (c *Context) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Context) Close() error {
	return c.conn.Close()
}

func (c *Context) HTML(code int, filename string, obj interface{}) {
	c.writeStatusCode(code)
	c.setContentType(htmlContentType)
	c.WriteHeaders()
	c.renderHTML(filename, obj)
}

func (c *Context) JSON(code int, obj interface{}) {
	c.writeStatusCode(code)
	c.setContentType(jsonContentType)
	c.WriteHeaders()
	c.WriteJSON(obj)
}

func (c *Context) WriteString(s string) {
	c.writeStatusCode(200)
	c.setContentType(htmlContentType)
	c.WriteHeaders()
	c.conn.Write([]byte(s))
}

func (c *Context) writeStatusCode(code int) {
	c.conn.Write([]byte("HTTP/1.1 " + strconv.Itoa(code) + " " + ReasonPhrase[code] + "\r\n"))
}

func (c *Context) WriteHeaders() {
	if !c.isHeaderWrite {
		for k, v := range c.Headers {
			c.conn.Write([]byte(k + ": " + v[0] + "\r\n"))
		}
		c.conn.Write([]byte{'\r', '\n'})
	}
	c.isHeaderWrite = true
}

func (c *Context) setContentType(val string) {
	if c.Headers.Get("Content-Type") == "" {
		c.Headers.Set("Content-Type", val)
	}
}

func (c *Context) renderHTML(filename string, obj interface{}) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("render html file error: ", err)
		return
	}
	c.conn.Write(f)
}

func (c *Context) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}

func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value, ok := c.GetQuery(key); ok {
		return value
	}
	return defaultValue
}

func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) QueryArray(key string) []string {
	values, _ := c.GetQueryArray(key)
	return values
}

func (c *Context) GetQueryArray(key string) ([]string, bool) {
	c.getQueryCache()
	if values, ok := c.queryCache[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

func (c *Context) getQueryCache() {
	if c.queryCache == nil {
		c.queryCache = c.Request.URL.Query()
	}
}

func (c *Context) PostForm(key string) string {
	value, _ := c.GetPostForm(key)
	return value
}

func (c *Context) DefaultPostForm(key, defaultValue string) string {
	if value, ok := c.GetPostForm(key); ok {
		return value
	}
	return defaultValue
}

func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) PostFormArray(key string) []string {
	values, _ := c.GetPostFormArray(key)
	return values
}

func (c *Context) getFormCache() {
	if c.formCache == nil {
		c.formCache = make(url.Values)
		req := c.Request
		if err := req.ParseFrom(); err != nil {
			log.Error(err)
		}
		c.formCache = req.PostForm
	}
}

// GetPostFormArray returns a slice of strings for a given form key, plus
// a boolean value whether at least one value exists for the given key.
func (c *Context) GetPostFormArray(key string) ([]string, bool) {
	c.getFormCache()
	if values := c.formCache[key]; len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

func (c *Context) NotFound() {
	c.writeStatusCode(404)
	c.setContentType(htmlContentType)
	c.WriteHeaders()
	c.conn.Write([]byte("404 not found"))
}

func (c *Context) WriteJSON(obj interface{}) error {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(jsonBytes)
	return err
}
