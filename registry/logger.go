package registry

type Logger interface {
	Info(...any)
	Debug(...any)
	Warn(...any)
	Error(...any)
	Fatal(...any)
}
