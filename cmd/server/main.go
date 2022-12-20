package main

import (
	"log"

	"github.com/kaikourok/lunchtote-backend/infrastructure/config"
	"github.com/kaikourok/lunchtote-backend/infrastructure/database"
	"github.com/kaikourok/lunchtote-backend/infrastructure/email"
	"github.com/kaikourok/lunchtote-backend/infrastructure/logger"
	"github.com/kaikourok/lunchtote-backend/infrastructure/notificator"
	"github.com/kaikourok/lunchtote-backend/registry"

	_ "github.com/lib/pq"
)

type serverRegisty struct {
	repository  *database.PostgresRepository
	config      *config.Config
	logger      *logger.Logger
	notificator *notificator.Notificator
	email       *email.MailSender
}

func (s serverRegisty) GetRepository() registry.Repository {
	return s.repository
}

func (s serverRegisty) GetConfig() registry.Config {
	return s.config
}

func (s serverRegisty) GetLogger() registry.Logger {
	return s.logger
}

func (s serverRegisty) GetNotificator() registry.Notificator {
	return s.notificator
}

func (s serverRegisty) GetEmail() registry.Email {
	return s.email
}

func main() {
	config, err := config.NewConfig("development")
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.NewLogger(false)
	if err != nil {
		log.Fatal(err)
	}

	notificator := notificator.NewNotificator(
		config.GetString("webhook.name"),
		config.GetString("webhook.avatar-url"),
	)

	email := email.NewMailSender(&email.MailConfig{
		Host:           config.GetString("email.smtphost"),
		Port:           config.GetInt("email.smtpport"),
		Address:        config.GetString("email.address"),
		Password:       config.GetString("email.password"),
		SenderName:     config.GetString("email.name"),
		ConnectTimeout: config.GetInt("email.connect-timeout"),
		SendTimeout:    config.GetInt("email.send-timeout"),
		DkimPrivateKey: config.GetString("email.dkim-private-key"),
		DkimSelector:   config.GetString("email.dkim-selector"),
	})

	repository, err := database.NewRepository(&database.DataSource{
		Host:     config.GetString("database.host"),
		Port:     config.GetInt("database.port"),
		User:     config.GetString("database.username"),
		Database: config.GetString("database.name"),
		Password: config.GetString("database.password"),
	})
	if err != nil {
		log.Fatal(err)
	}

	registry := serverRegisty{repository, config, logger, notificator, email}

	router := NewRouter(registry)
	router.Run(":" + config.GetString("general.port"))
}
