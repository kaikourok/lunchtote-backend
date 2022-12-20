package notificator

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Notificator struct {
	DisplayName string
	AvatarUrl   string
}

func NewNotificator(displayName, avatarUrl string) *Notificator {
	return &Notificator{displayName, avatarUrl}
}

type discordWebhook struct {
	UserName  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
	Content   string `json:"content"`
}

func (n *Notificator) SendWebhook(url, message string) {
	jsonString, err := json.Marshal(discordWebhook{
		UserName:  n.DisplayName,
		AvatarUrl: n.AvatarUrl,
		Content:   message,
	})
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	client.Do(req)
}
