package myhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
	"strings"
	"website/log"
)

type Context struct {
	conn           Conn
	body           []byte
	Headers        Header
	Request        *Request
	isRange        bool
	isPartialRange bool
	isHeaderWrite  bool
	rangeL         int
	rangeR         int
	ContentLength  int
	StatusCode     int
	queryCache     url.Values
	formCache      url.Values
}

const (
	htmlContentType  = "text/html; charset=utf-8"
	jsonContentType  = "application/json; charset=utf-8"
	videoContentType = "application/mp4"
)

var (
	ErrRangeNaN   = errors.New("header range not number")
	ErrRangeError = errors.New("header range negative or overlap or too big")
)

func NewContext(c Conn, req *Request) *Context {
	return &Context{c, make([]byte, 0), make(map[string][]string), req, false, false, false, 0, 0, 0, 0, nil, nil}
}

func (c *Context) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Context) Close() error {
	return c.conn.Close()
}

func (c *Context) HTML(code int, templateName string, templates []string, obj interface{}) {
	c.StatusCode = code
	c.setContentType(htmlContentType)
	c.renderHTML(templateName, templates, obj)
	c.WriteHeaders()
	c.WriteBody()
}

func (c *Context) JSON(code int, obj interface{}) {
	c.StatusCode = code
	c.setContentType(jsonContentType)
	c.WriteJSON(obj)
	c.WriteHeaders()
	c.WriteBody()
}

func (c *Context) VIDEO(filepath string) {
	c.StatusCode = 200
	c.setContentType(videoContentType)
	c.writeFile(filepath)
	c.handleRange()
	c.WriteHeaders()
	c.WriteBody()
}

func (c *Context) WriteString(s string) {
	c.StatusCode = 200
	c.setContentType(htmlContentType)
	c.WriteHeaders()
	c.conn.Write([]byte(s))
}

func (c *Context) writeStatusCode() {
	c.conn.Write([]byte("HTTP/1.1 " + strconv.Itoa(c.StatusCode) + " " + ReasonPhrase[c.StatusCode] + "\r\n"))
}

func (c *Context) WriteHeaders() {
	if !c.isHeaderWrite {
		c.writeStatusCode()
		if c.ContentLength == 0 {
			c.ContentLength = len(c.body)
		}
		c.Headers.Set("Content-Length", strconv.Itoa(c.ContentLength))
		for k, v := range c.Headers {
			c.conn.Write([]byte(k + ": " + v[0] + "\r\n"))
		}
		c.conn.Write([]byte{'\r', '\n'})
	}
	c.isHeaderWrite = true
}

func (c *Context) WriteBody() {
	c.conn.Write(c.body)
}

func (c *Context) setContentType(val string) {
	if c.Headers.Get("Content-Type") == "" {
		c.Headers.Set("Content-Type", val)
	}
}

func (c *Context) renderHTML(templateName string, templates []string, obj interface{}) {
	tpl := template.Must(template.ParseFiles(templates...))
	var t bytes.Buffer
	err := tpl.ExecuteTemplate(&t, templateName, obj)
	if err != nil {
		log.Error("render html execute error: ", err)
		return
	}
	c.body = t.Bytes()

}

func (c *Context) writeFile(filepath string) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Error(err)
		return
	}
	c.body = data
}

func (c *Context) RequestIP() net.IP {
	return c.conn.Addr()
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
		if err := req.ParseForm(); err != nil {
			log.Error(err)
		}
		cleanForm := make(url.Values)
		for k, v := range req.PostForm {
			var newV []string
			for _, vv := range v {
				newV = append(newV, cleanNullByte(vv))
			}
			cleanForm[k] = newV
		}
		c.formCache = cleanForm
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
	c.StatusCode = 404
	c.setContentType(htmlContentType)
	c.WriteHeaders()
	c.conn.Write([]byte("404 Not Found"))
}

func (c *Context) Forbidden() {
	c.StatusCode = 403
	c.setContentType(htmlContentType)
	c.WriteHeaders()
	c.conn.Write([]byte("403 Forbidden"))
}

func (c *Context) InternalError() {
	c.StatusCode = 500
	c.setContentType(htmlContentType)
	c.WriteHeaders()
	c.conn.Write([]byte("500 Internal Error"))
}

func (c *Context) WriteJSON(obj interface{}) error {
	if obj == nil {
		return nil
	}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	c.body = jsonBytes
	return nil
}

func (c *Context) SetCookie(key, val string) {

}

func (c *Context) Cookie(key string) string {
	cookie, ok := c.Request.Headers["Cookie"]
	if !ok {
		return ""
	}
	s := strings.Split(cookie[0], "; ")
	for i := range s {
		ss := strings.SplitN(s[i], "=", 2)
		k, v := ss[0], ss[1]
		if len(k) == 0 {
			continue
		}
		if k == key {
			return v
		}
	}
	return ""
}

func (c *Context) handleRange() error {
	c.Headers.Set("Content-Transfer-Encoding", "binary")
	if v, ok := c.Request.Headers["Range"]; ok && len(v) == 1 && strings.HasPrefix(v[0], "bytes=") {
		req := v[0][6:]
		sp := strings.SplitN(req, "-", 2)
		startSizeString, endSizeString := sp[0], sp[1]
		startSize, err := strconv.Atoi(startSizeString)
		if err != nil {
			return ErrRangeNaN
		}
		endSize := len(c.body) - 1
		if len(endSizeString) > 0 {
			endSize, err = strconv.Atoi(endSizeString)
			if err != nil {
				return ErrRangeNaN
			}
		} else {
			endSizeString = strconv.Itoa(endSize)
		}
		if startSize < 0 || startSize > endSize {
			return ErrRangeError
		}
		c.Headers.Set("Accept-Ranges", "bytes")
		size := strconv.Itoa(len(c.body))
		c.Headers.Set("Content-Range", "bytes "+startSizeString+"-"+endSizeString+"/"+size)
		c.body = c.body[startSize : endSize+1]
		c.ContentLength = endSize - startSize + 1
		c.StatusCode = 206
	}
	return nil
}

func cleanNullByte(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != 0 {
			return s[:i+1]
		}
	}
	return s
}
