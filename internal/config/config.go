package config

import (
	"fmt"
	"middleman/internal/rpc"
	"middleman/internal/web"
	"os"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
)

type (
	PGConfig struct {
		Conn string
	}

	NatsConfig struct {
		URL    string `required:"true"`
		Stream string `default:"middleman"`
	}

	OtelConfig struct {
		ServiceName      string `envconfig:"SERVICE_NAME" default:"middleman"`
		ExporterEndpoint string `envconfig:"EXPORTER_OTLP_ENDPOINT" default:"http://collector:4317"`
	}

	AppConfig struct {
		Environment     string
		LogLevel        string `envconfig:"LOG_LEVEL" default:"DEBUG"`
		PG              PGConfig
		Nats            NatsConfig
		Rpc             rpc.RpcConfig
		Web             web.WebConfig
		Otel            OtelConfig
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`

		JWTSecret        string        `envconfig:"JWT_SECRET_KEY" required:"true"`
		JWTIssuer        string        `envconfig:"JWT_ISSUER" required:"true"`
		JWTAudience      string        `envconfig:"JWT_AUDIENCE" required:"true"`
		JWTTokenExpiry   time.Duration `envconfig:"JWT_TOKEN_EXPIRY" default:"1h"`
		JWTRefreshExpiry time.Duration `envconfig:"JWT_REFRESH_EXPIRY" default:"168h"`

		CookieDomain string `envconfig:"COOKIE_DOMAIN" required:"true"`
		CookiePath   string `envconfig:"COOKIE_PATH" default:"/"`
		CookieName   string `envconfig:"COOKIE_NAME" default:"auth_token"`

		APIKey string `envconfig:"API_KEY" required:"true"`
	}
)

// detectServiceType detects which service we're running based on executable path or environment
func detectServiceType() string {
	if serviceName := os.Getenv("SERVICE_NAME"); serviceName != "" {
		return serviceName
	}

	// Fallback: try to detect from executable path or working directory
	if wd, err := os.Getwd(); err == nil {
		if strings.Contains(wd, "vectors") {
			return "vectors"
		}
	}

	return "unknown"
}

func InitConfig() (cfg AppConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}

	err = envconfig.Process("", &cfg)
	if err != nil {
		return
	}

	// Post-process validation: check if PG_CONN is required for this service
	serviceType := detectServiceType()
	if serviceType != "vectors" && cfg.PG.Conn == "" {
		return cfg, fmt.Errorf("required key PG_CONN missing value (required for %s service)", serviceType)
	}

	return
}
