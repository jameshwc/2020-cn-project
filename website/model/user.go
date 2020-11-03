package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id"`
	Name     string        `json:"username"`
	Password string
}

type Auth struct {
	ID                  bson.ObjectId `bson:"_id"`
	UID                 bson.ObjectId `bson:"uid"`
	SessionID           uuid.UUID     `bson:"session_id"`
	CreatedAt           time.Time     `bson:"created_at"`
	ExpiredAfterSeconds uint          `bson:"expireAfterSeconds"`
}

func NewUser(name, password string) *User {
	return &User{bson.NewObjectId(), name, encrypt(password)}
}

// TODO: encrypt the password
func encrypt(password string) string {
	return password
}

func (u *User) Validate(password string) bool {
	return encrypt(password) == u.Password
}

func NewAuth(uid bson.ObjectId) *Auth {
	return &Auth{bson.NewObjectId(), uid, uuid.NewV4(), time.Now(), 3600}
}
