package model

import "time"

type MailBase struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Read      bool      `json:"read"`
}

type ReceivedMail struct {
	MailBase
	Sender struct {
		Id   *int   `json:"id"`
		Name string `json:"name"`
	} `json:"sender"`
}

type SentMail struct {
	MailBase
	Receiver struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"receiver"`
}
