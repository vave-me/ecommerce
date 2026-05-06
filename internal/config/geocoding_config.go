// payments/internal/config/payments_config.go
package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
)

// GeocodingConfig extends the shared AppConfig with Stripe-specific configurations
type GeocodingConfig struct {
	GoogleAPIKey string `envconfig:"GEOCODING_API_KEY" required:"true"`
	GoogleSecret string `envconfig:"GEOCODING_TOKEN_SECRET" required:"true"`
}

// InitGeocodingConfig initializes the GeocodingConfig by loading shared and Stripe-specific configurations
func InitGeocodingConfig() (paymentCfg GeocodingConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &paymentCfg)
	return
}
