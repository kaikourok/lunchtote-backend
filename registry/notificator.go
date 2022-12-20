package registry

type Notificator interface {
	SendWebhook(url, message string)
}
