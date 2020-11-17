package view

import (
	"website/myhttp"
	viewuser "website/view/user"
)

func Audio(c *myhttp.Context) {
	c.MP3("static/audio/wake_me_up_320k.mp3")
}

func AudioDemo(c *myhttp.Context) {
	auth, username := viewuser.CheckLogin(c.Cookie("sessionid"))

	c.HTML(200, "video", []string{"static/audio.html", "static/header.html", "static/footer.html"}, map[string]interface{}{
		"PageTitle": "Audio", "Auth": auth, "Username": username,
	})
}
