package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
)

// UsersConfig extends the shared AppConfig with Stripe-specific configurations
type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST"`
	Port     string `envconfig:"REDIS_PORT"`
	Password string `envconfig:"REDIS_PASSWORD" `
	Database int    `envconfig:"REDIS_DATABASE"`
}

// InitCommentsConfig initializes the CommentsConfig by loading shared and Stripe-specific configurations
func InitRedisConfig() (redisConfig RedisConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &redisConfig)
	return
}
