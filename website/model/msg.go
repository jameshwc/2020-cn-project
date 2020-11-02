package model

import "time"

type Message struct {
	UserName   string    `json:"user_name"`
	Name       string    `json:"name"`
	Content    string    `json:"content"`
	SubmitTime time.Time `json:"submit_time"`
}

func NewMessage(userName, name, content string) *Message {
	return &Message{cleanNullByte(userName), cleanNullByte(name), cleanNullByte(content), time.Now()}
}

func cleanNullByte(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != 0 {
			return s[:i+1]
		}
	}
	return s
}
