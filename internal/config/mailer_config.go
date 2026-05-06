package config

import (
	"fmt"
	"gopkg.in/mail.v2"
	"os"
)

type MailerConfig struct {
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	DefaultSender string
	ServerAddress string
}

func LoadMailerConfig() MailerConfig {
	fmt.Println("Loading Mailer environment variables...")

	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp-relay.brevo.com"
		fmt.Println("SMTP_HOST not set, using default:", smtpHost)
	}

	smtpPort := 587 // Brevo default
	if os.Getenv("SMTP_PORT") != "" {
		fmt.Sscanf(os.Getenv("SMTP_PORT"), "%d", &smtpPort)
	}

	smtpUsername := os.Getenv("SMTP_USERNAME")
	if smtpUsername == "" {
		smtpUsername = "redacted-email@example.com"
		fmt.Println("SMTP_USERNAME not set, using default:", smtpUsername)
	}

	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		fmt.Println("WARNING: SMTP_PASSWORD not set!")
	}

	defaultSender := os.Getenv("DEFAULT_SMTP_SENDER")
	if defaultSender == "" {
		defaultSender = "redacted-email@example.com"
		fmt.Println("DEFAULT_SMTP_SENDER not set, using default:", defaultSender)
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = ":9000"
		fmt.Println("SERVER_ADDRESS not set, using default:", serverAddress)
	}

	return MailerConfig{
		SMTPHost:      smtpHost,
		SMTPPort:      smtpPort,
		SMTPUsername:  smtpUsername,
		SMTPPassword:  smtpPassword,
		DefaultSender: defaultSender,
		ServerAddress: serverAddress,
	}
}

func NewSMTPDialer(cfg MailerConfig) *mail.Dialer {
	d := mail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	return d
}
