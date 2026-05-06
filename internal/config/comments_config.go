package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
)

// CommentsConfig extends the shared AppConfig with Stripe-specific configurations
type CommentsConfig struct {
	WEBSOCKET        string `envconfig:"NATS_WEBSOCKET"`
	StreamName       string `envconfig:"COMMENTS_STREAM"`
	WebSocketSubject string `envconfig:"COMMENTS_WEBSOCKET_SUBJECT"`
	RedisHost        string `envconfig:"REDIS_HOST"`
	RedisPort        string `envconfig:"REDIS_PORT"`
	RedisPassword    string `envconfig:"REDIS_PASSWORD"`
	RedisDatabase    int    `envconfig:"REDIS_DATABASE" default:"2"`
}

// InitCommentsConfig initializes the CommentsConfig by loading shared and Stripe-specific configurations
func InitCommentsConfig() (commentsConfig CommentsConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &commentsConfig)
	return
}
