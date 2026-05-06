// payments/internal/config/payments_config.go
package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
)

// PaymentsConfig extends the shared AppConfig with Stripe-specific configurations
type MerchantConfig struct {
	MerchantId             uint64 `envconfig:"MERCHANT_ID" required:"true"`
	ServiceAccountJSONPath string `envconfig:"SERVICE_ACCOUNT_JSON_PATH" default:""`
	DevelopmentMode        bool   `envconfig:"MERCHANT_DEVELOPMENT_MODE" default:"true"`
}

// InitPaymentsConfig initializes the PaymentsConfig by loading shared and Stripe-specific configurations
func InitMerchantConfig() (merchantConfig MerchantConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &merchantConfig)
	return
}
