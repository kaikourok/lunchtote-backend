package registry

type Email interface {
	SendEmail(to, subject, message string) error
}
