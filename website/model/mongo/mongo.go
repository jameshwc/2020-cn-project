package mongo

import (
	"log"
	"website/conf"

	"gopkg.in/mgo.v2"
)

var db *mgo.Database

func Setup() {
	session, err := mgo.Dial(conf.DBconfig.Host)
	if err != nil {
		log.Fatal("mongo: ", err)
	}
	session.SetMode(mgo.Monotonic, true)
	db = session.DB(conf.DBconfig.Name)
	if err := db.Session.Ping(); err != nil {
		log.Fatal("mongo: ", err)
	}
}
