package system

import (
	"context"
	"database/sql"
	ai2 "middleman/internal/ai"
	"middleman/internal/auth"
	"middleman/internal/config"
	"middleman/internal/geo"
	"middleman/internal/merchant"
	oidcclient "middleman/internal/oid"
	"middleman/internal/stripe"
	"middleman/internal/waiter"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Service interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
}

type Module interface {
	Startup(context.Context, Service) error
}

type PaymentService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	Stripe() stripe.StripeClient
	PaymentConfig() config.PaymentsConfig
}

type MerchantService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	Merchant() merchant.MerchantCenterClient
	MerchantConfig() config.MerchantConfig
}
type UsersService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	WebGoogleOIDC() *oidcclient.GoogleOIDCClient
	InitWebGoogleOIDC() error
	MobileGoogleOIDC() *oidcclient.GoogleOIDCClient
	InitMobileGoogleOIDC() error
	UsersConfig() config.UsersConfig
}
type ShippingService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	DHLClient() config.DHLClient
	ShippingConfig() config.ShippingConfig
}

type MessengerService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	MessengerConfig() config.MessengerConfig
}
type CommentsService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	RedisPool() *redis.Pool
	CommentsConfig() config.CommentsConfig
}

type ActivityService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	RedisPoolActivity() *redis.Pool
	ActivityConfig() config.ActivityConfig
}

type MediaService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	MediaConfig() config.MediaConfig
}

type MailerService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	MailerConfig() config.MailerConfig
}
type SearchService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	Redis() *redis.Pool
	Redisearch() *redisearch.Client
	SearchConfig() config.SearchConfig
}
type GeocodingService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	Geocode() geo.GoogleGeocodingClient
	GeocodingConfig() config.GeocodingConfig
}
type MetricsService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	Redis() *redis.Pool
	MetricsConfig() config.MetricsConfig
}
type AssistantsService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	Auth() *auth.Auth
	AssistantsConfig() config.AssistantsConfig
	AnthropicClient() *ai2.AnthropicClient
	DeepSeekClient() *ai2.DeepSeekClient
	OpenAiClient() *ai2.OpenAIClient
}
type VectorsService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	Auth() *auth.Auth
	VectorsConfig() config.VectorsConfig
	AnthropicClient() *ai2.AnthropicClient
	DeepSeekClient() *ai2.DeepSeekClient
	OpenAiClient() *ai2.OpenAIClient
	Redis() *redis.Pool
	Redisearch() *redisearch.Client
}

type SchedulerService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	RedisPoolScheduler() *redis.Pool
	SchedulerConfig() config.SchedulerConfig
}
type ManagersService interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
	Auth() *auth.Auth
	ManagersConfig() config.ManagersConfig
	AnthropicClient() *ai2.AnthropicClient
	DeepSeekClient() *ai2.DeepSeekClient
	OpenAiClient() *ai2.OpenAIClient
}
