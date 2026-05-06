package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
)

// MetricsConfig extends the shared AppConfig with redis-specific configurations
type MetricsConfig struct {
	IndexName     string `envconfig:"REDISEARCH_INDEX_NAME" default:"metrics_index`
	RedisHost     string `envconfig:"REDIS_HOST"`
	RedisPort     string `envconfig:"REDIS_PORT"`
	RedisPassword string `envconfig:"REDIS_PASSWORD"`
	RedisDatabase int    `envconfig:"REDIS_DATABASE"`
}

// InitMetricsConfig initializes the MetricsConfig by loading shared and Stripe-specific configurations
func InitMetricsConfig() (metricsConfig MetricsConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &metricsConfig)
	return
}
