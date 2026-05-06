package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
)

// SchedulerConfig extends the shared AppConfig with scheduler-specific configurations
type SchedulerConfig struct {
	RedisHost     string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort     string `envconfig:"REDIS_PORT" default:"6379"`
	RedisPassword string `envconfig:"REDIS_PASSWORD"`
	RedisDB       int    `envconfig:"REDIS_DB" default:"0"`
	AssistantID   string `envconfig:"SCHEDULER_ASSISTANT_ID" default:"scheduler-assistant"`
}

// InitSchedulerConfig initializes the SchedulerConfig by loading shared and scheduler-specific configurations
func InitSchedulerConfig() (schedulerConfig SchedulerConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &schedulerConfig)
	return
}