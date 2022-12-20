package registry

type Registry interface {
	GetRepository() Repository
	GetConfig() Config
	GetLogger() Logger
	GetNotificator() Notificator
	GetEmail() Email
}
