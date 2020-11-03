package view

import (
	"fmt"
	"website/log"
	"website/model"
	"website/model/mongo"
	"website/myhttp"
	viewuser "website/view/user"
)

func AddMessage(c *myhttp.Context) {
	fmt.Println(c.PostForm("content"), len(c.PostForm("content")))
	msg := model.NewMessage("james", c.PostForm("name"), c.PostForm("content"))
	c.Headers.Set("Location", "/messages")
	err := mongo.AddMessage(msg)
	if err != nil {
		c.InternalError()
		log.InfoWithSource(err)
		return
	}
	c.WriteString("submit successfully")
}

func ShowMessage(c *myhttp.Context) {
	auth, username := viewuser.CheckLogin(c.Cookie("sessionid"))

	msg, err := mongo.GetMessageAll()
	if err != nil {
		c.InternalError()
		log.InfoWithSource(err)
		return
	}
	c.HTML(200, "message", []string{"static/message.html", "static/header.html", "static/footer.html"}, map[string]interface{}{
		"PageTitle": "Message Board", "messages": msg, "Auth": auth, "Username": username,
	})
}
