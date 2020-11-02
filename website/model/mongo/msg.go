package mongo

import (
	"website/model"
)

const msgTableName = "messages"

func AddMessage(msg *model.Message) error {
	c := db.C(msgTableName)
	err := c.Insert(msg)
	if err != nil {
		return err
	}
	return nil
}

func GetMessageAll() (msg []model.Message, err error) {
	c := db.C(msgTableName)
	err = c.Find(nil).All(&msg)
	return
}
