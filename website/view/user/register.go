package viewuser

import (
	"website/model"
	"website/model/mongo"
	"website/myhttp"
)

func RegisterForm(c *myhttp.Context) {
	auth, username := CheckLogin(c.Cookie("sessionid"))

	c.HTML(200, "user/register_form", []string{"static/header.html", "static/register_form.html", "static/footer.html"}, map[string]interface{}{
		"PageTitle": "Register", "Auth": auth, "Username": username,
	})
}

func Register(c *myhttp.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if len(username) == 0 || len(password) == 0 {
		c.Forbidden()
		return
	}
	user := model.NewUser(username, password)
	mongo.AddUser(user)
	c.WriteString("register successfully")
}
