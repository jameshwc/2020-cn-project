package viewuser

import (
	"website/model/mongo"
	"website/myhttp"
)

func LoginForm(c *myhttp.Context) {
	auth, username := CheckLogin(c.Cookie("sessionid"))

	c.HTML(200, "user/login_form", []string{"static/header.html", "static/login_form.html", "static/footer.html"}, map[string]interface{}{
		"PageTitle": "Login", "Auth": auth, "Username": username,
	})
}

func Login(c *myhttp.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if len(username) == 0 || len(password) == 0 {
		c.Forbidden()
		return
	}
	user, sessionID, err := mongo.Login(username, password)
	if err != nil {
		c.InternalError()
		return
	} else if user == nil {
		c.Forbidden()
		return
	}
	c.Headers.Set("Set-Cookie", "sessionid="+sessionID.String()+"; ")
	c.WriteString("Login successfully")
}
