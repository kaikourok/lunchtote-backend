package main

import (
	"fmt"
	"log"

	"github.com/kaikourok/lunchtote-backend/infrastructure/config"
	"github.com/kaikourok/lunchtote-backend/infrastructure/database"
	"github.com/kaikourok/lunchtote-backend/infrastructure/email"
	"github.com/kaikourok/lunchtote-backend/infrastructure/logger"
	"github.com/kaikourok/lunchtote-backend/infrastructure/notificator"
	"github.com/kaikourok/lunchtote-backend/registry"
	"github.com/kaikourok/lunchtote-backend/usecase/control"
	"github.com/spf13/cobra"
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

var rootCmd = &cobra.Command{
	Use:   "Lunchtote CLI",
	Short: "Control lunchtote",
	Long:  "Control lunchtote",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("requires subcommands.")
	},
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

	rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "システムを初期化します。",
		Long:  "システムを初期化します。",
		Run: func(cmd *cobra.Command, args []string) {
			usecase := control.NewControlUsecase(registry)

			fmt.Println("初期化を開始します。")
			err := usecase.Initialize()
			if err != nil {
				fmt.Println("エラーが発生しました。強制終了します。")
				log.Fatal(err)
			}

			fmt.Println("初期化が完了しました。")
		},
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
