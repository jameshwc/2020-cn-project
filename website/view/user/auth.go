package viewuser

import "website/model/mongo"

func CheckLogin(sessionID string) (auth bool, username string) {
	u, err := mongo.Auth(sessionID)
	if err != nil {
		return
	} else if u != nil {
		auth = true
		username = u.Name
	}
	return
}
