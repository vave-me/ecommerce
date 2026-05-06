// payments/internal/config/payments_config.go
package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
)

// PaymentsConfig extends the shared AppConfig with Stripe-specific configurations
type PaymentsConfig struct {
	StripeAPIKey        string `envconfig:"STRIPE_KEY" required:"true"`
	StripeWebhookSecret string `envconfig:"STRIPE_SECRET" required:"true"`
	WebhookSecret       string `envconfig:"STRIPE_WEBHOOK_SECRET" required:"true"`
}

// InitPaymentsConfig initializes the PaymentsConfig by loading shared and Stripe-specific configurations
func InitPaymentsConfig() (paymentCfg PaymentsConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &paymentCfg)
	return
}
