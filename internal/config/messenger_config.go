package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
)

// MessengerConfig extends the shared AppConfig with Stripe-specific configurations
type MessengerConfig struct {
	WEBSOCKET        string `envconfig:"NATS_WEBSOCKET"`
	StreamName       string `envconfig:"MESSENGER_STREAM"`
	WebSocketSubject string `envconfig:"MESSENGER_WEBSOCKET_SUBJECT"`
}

// InitMessengerConfig initializes the MessengerConfig by loading shared and Stripe-specific configurations
func InitMessengerConfig() (messengerConfig MessengerConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &messengerConfig)
	return
}
