package view

import (
	"website/myhttp"
	viewuser "website/view/user"
)

func Homepage(c *myhttp.Context) {
	auth, username := viewuser.CheckLogin(c.Cookie("sessionid"))

	c.HTML(200, "index", []string{"static/index.html", "static/header.html", "static/footer.html"}, map[string]interface{}{
		"PageTitle": "Home", "Auth": auth, "Username": username,
	})
}
