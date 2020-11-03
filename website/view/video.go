package view

import (
	"website/myhttp"
	viewuser "website/view/user"
)

func Video(c *myhttp.Context) {
	c.VIDEO("static/video/video2.mp4")
}

func VideoDemo(c *myhttp.Context) {
	auth, username := viewuser.CheckLogin(c.Cookie("sessionid"))

	c.HTML(200, "video", []string{"static/video.html", "static/header.html", "static/footer.html"}, map[string]interface{}{
		"PageTitle": "Video", "Auth": auth, "Username": username,
	})
}
