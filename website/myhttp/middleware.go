package myhttp

import (
	"website/log"
)

func beforeServe(c *Context) {
	log.TraceIP(c.RequestIP().String(), c.Request.Method, c.Request.URL.Path)
}

func afterServe(c *Context) {

}
