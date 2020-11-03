package viewuser

import (
	"website/model/mongo"
	"website/myhttp"
)

func Logout(c *myhttp.Context) {
	sessionID := c.Cookie("sessionid")
	auth, _ := CheckLogin(sessionID)
	if auth {
		err := mongo.DeleteAuth(sessionID)
		if err != nil {
			c.InternalError()
			return
		}
		c.Headers.Set("Set-Cookie", "sessionid=; ")
		c.HTML(200, "user/logout", []string{"static/header.html", "static/logout.html", "static/footer.html"}, map[string]interface{}{
			"PageTitle": "Logout", "Auth": false, "Username": "", "Logout": true,
		})
	} else {
		c.HTML(403, "user/logout", []string{"static/header.html", "static/logout.html", "static/footer.html"}, map[string]interface{}{
			"PageTitle": "Logout", "Auth": false, "Username": "", "Logout": false,
		})
	}
}
