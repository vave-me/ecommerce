package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
)

// ActivityConfig extends the shared AppConfig with Stripe-specific configurations
type ActivityConfig struct {
	RedisHost     string `envconfig:"REDIS_HOST"`
	RedisPort     string `envconfig:"REDIS_PORT"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" `
	RedisDatabase int    `envconfig:"REDIS_DATABASE" default:"1"`
}

// InitActivityConfig initializes the ActivityConfig by loading shared and Stripe-specific configurations
func InitActivityConfig() (activityConfig ActivityConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &activityConfig)
	return
}
