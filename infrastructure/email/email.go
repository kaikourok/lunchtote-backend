package email

import (
	"time"

	"github.com/toorop/go-dkim"
	mail "github.com/xhit/go-simple-mail/v2"
)

type MailConfig struct {
	Host           string
	Port           int
	Address        string
	Password       string
	SenderName     string
	ConnectTimeout int
	SendTimeout    int
	DkimPrivateKey string
	DkimSelector   string
}

type MailSender struct {
	MailConfig
}

func (s *MailSender) SendEmail(to, subject, message string) error {
	server := mail.NewSMTPClient()

	server.Host = s.Host
	server.Port = s.Port
	server.Username = s.Address
	server.Password = s.Password
	server.Encryption = mail.EncryptionSTARTTLS

	server.KeepAlive = false
	server.ConnectTimeout = time.Duration(s.ConnectTimeout) * time.Second
	server.SendTimeout = time.Duration(s.SendTimeout) * time.Second

	client, err := server.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	email := mail.NewMSG()
	email.SetFrom(s.SenderName + " <" + s.Address + ">").
		AddTo(to).
		SetSubject(subject)

	email.SetBody(mail.TextPlain, message)

	if s.DkimPrivateKey != "" {
		options := dkim.NewSigOptions()
		options.PrivateKey = []byte(s.DkimPrivateKey)
		options.Domain = s.Host
		options.Selector = s.DkimSelector
		options.SignatureExpireIn = 3600
		options.Headers = []string{"from", "date", "mime-version", "received", "received"}
		options.AddSignatureTimestamp = true
		options.Canonicalization = "relaxed/relaxed"

		email.SetDkim(options)
	}

	if email.Error != nil {
		return email.Error
	}

	err = email.Send(client)

	if err != nil {
		return err
	}

	return nil
}

func NewMailSender(config *MailConfig) *MailSender {
	return &MailSender{
		*config,
	}
}
