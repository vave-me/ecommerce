package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
)

// SearchConfig extends the shared AppConfig with redis-specific configurations
type SearchConfig struct {
	IndexName string `envconfig:"REDISEARCH_INDEX_NAME" default:"search_index`
	Host      string `envconfig:"REDIS_HOST" `
	Port      string `envconfig:"REDIS_PORT"`
	Password  string `envconfig:"REDIS_PASSWORD"`
}

// InitSearchConfig initializes the SearchConfig by loading shared and Stripe-specific configurations
func InitSearchConfig() (searchConfig SearchConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &searchConfig)
	return
}
