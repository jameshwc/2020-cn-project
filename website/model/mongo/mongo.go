package mongo

import (
	"log"
	"website/conf"

	"gopkg.in/mgo.v2"
)

var db *mgo.Database
var uC *mgo.Collection
var aC *mgo.Collection
var mC *mgo.Collection

func Setup() {
	session, err := mgo.Dial(conf.DBconfig.Host)
	if err != nil {
		log.Fatal("mongo: ", err)
	}
	session.SetMode(mgo.Monotonic, true)
	db = session.DB(conf.DBconfig.Name)
	if err := db.Session.Ping(); err != nil {
		log.Fatal("mongo: ", err, " ", conf.DBconfig.Host)
	}
	uC = db.C(userCollectionName)
	aC = db.C(authCollectionName)
	mC = db.C(msgCollectionName)
}
