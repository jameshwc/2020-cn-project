package view

import (
	"fmt"
	"website/log"
	"website/model"
	"website/model/mongo"
	"website/myhttp"
)

func AddMessage(c *myhttp.Context) {
	fmt.Println(c.PostForm("content"), len(c.PostForm("content")))
	msg := model.NewMessage("james", c.PostForm("name"), c.PostForm("content"))
	c.Headers.Set("Location", "/messages")
	err := mongo.AddMessage(msg)
	if err != nil {
		c.HTML(500, "static/message_internal_error.html", nil)
		log.InfoWithSource(err)
		return
	}
	c.HTML(200, "static/message.html", nil)
}

func ShowMessage(c *myhttp.Context) {
	msg, err := mongo.GetMessageAll()
	if err != nil {
		c.HTML(500, "static/message_internal_error.html", nil)
		log.InfoWithSource(err)
		return
	}
	if len(msg) == 0 {
		c.HTML(500, "static/message_internal_error.html", nil)
		return
	}
	c.JSON(200, msg)
}
