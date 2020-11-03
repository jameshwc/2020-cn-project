package mongo

import (
	"website/model"
)

const msgCollectionName = "messages"

func AddMessage(msg *model.Message) error {
	err := mC.Insert(msg)
	if err != nil {
		return err
	}
	return nil
}

func GetMessageAll() (msg []model.Message, err error) {
	err = mC.Find(nil).All(&msg)
	return
}
