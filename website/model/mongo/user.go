package mongo

import (
	"fmt"
	"website/model"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	userCollectionName = "users"
	authCollectionName = "auth"
)

func AddUser(u *model.User) error {
	err := uC.Insert(u)
	if err != nil {
		return err
	}
	return nil
}

func Login(username, password string) (*model.User, uuid.UUID, error) {
	var u model.User
	err := uC.Find(bson.M{"name": username}).One(&u)

	if err != nil {
		return nil, uuid.Nil, err
	}
	if u.Validate(password) {
		fmt.Println("login successfully")
		auth := model.NewAuth(u.ID)
		aC.Insert(auth)
		return &u, auth.SessionID, nil
	}
	fmt.Println("Login Failed", u.Password, password, len(u.Password), len(password))
	return nil, uuid.Nil, nil
}

func Auth(sessionID string) (*model.User, error) {
	suuid := uuid.FromStringOrNil(sessionID)
	var a model.Auth
	err := aC.Find(bson.M{"session_id": suuid}).One(&a)
	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var u model.User
	err = uC.Find(bson.M{"_id": a.UID}).One(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func DeleteAuth(sessionID string) error {
	suuid := uuid.FromStringOrNil(sessionID)
	return aC.Remove(bson.M{"session_id": suuid})
}
