package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
)

// UsersConfig extends the shared AppConfig with Stripe-specific configurations
type WebhookConfig struct {
	WEBSOCKET        string `envconfig:"NATS_WEBSOCKET"`
	StreamName       string `envconfig:"COMMENTS_STREAM"`
	WebSocketSubject string `envconfig:"COMMENTS_WEBSOCKET_SUBJECT"`
}

// InitCommentsConfig initializes the CommentsConfig by loading shared and Stripe-specific configurations
func InitWebhookConfig() (webhookConfig WebhookConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &webhookConfig)
	return
}
