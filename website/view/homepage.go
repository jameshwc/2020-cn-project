package view

import "website/myhttp"

func Homepage(c *myhttp.Context) {
	c.HTML(200, "static/index.html", nil)
}
