// users/internal/config/users_config.go   ← (rename if you like)
// If this file really lives under payments/, keep the old path.
// ----------------------------------------------------------------
package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
)

// UsersConfig now only needs your Google OAuth Client ID.
type UsersConfig struct {
	WebGoogleOAuthClientID    string `envconfig:"WEB_GOOGLE_OAUTH_CLIENT_ID"`
	MobileGoogleOAuthClientID string `envconfig:"MOBILE_GOOGLE_OAUTH_CLIENT_ID"`
	Issuer                    string `envconfig:"GOOGLE_ISSUER"`
}

// InitUsersConfig loads .env (for dev/staging) and populates UsersConfig.
func InitUsersConfig() (usersCfg UsersConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &usersCfg)
	return
}
