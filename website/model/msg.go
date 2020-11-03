package model

import "time"

type Message struct {
	UserName   string    `json:"user_name"`
	Name       string    `json:"name"`
	Content    string    `json:"content"`
	SubmitTime time.Time `json:"submit_time"`
}

func NewMessage(userName, name, content string) *Message {
	return &Message{userName, name, content, time.Now()}
}
