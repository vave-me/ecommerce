package config

import (
	"time"
	
	"github.com/kelseyhightower/envconfig"
)

// Config holds all configuration for the streams service
type Config struct {
	// Server configuration
	HTTPPort string `envconfig:"HTTP_PORT" default:"8080"`
	GRPCPort string `envconfig:"GRPC_PORT" default:"9090"`
	
	// Environment
	Environment string `envconfig:"ENVIRONMENT" default:"production"`
	
	// Streaming configuration
	RTMP struct {
		Port           int    `envconfig:"RTMP_PORT" default:"1935"`
		MaxConnections int    `envconfig:"RTMP_MAX_CONNECTIONS" default:"1000"`
		AuthEnabled    bool   `envconfig:"RTMP_AUTH_ENABLED" default:"true"`
		AuthSecret     string `envconfig:"RTMP_AUTH_SECRET"`
	}
	
	SRT struct {
		Port        int    `envconfig:"SRT_PORT" default:"4578"`
		Passphrase  string `envconfig:"SRT_PASSPHRASE"`
		MaxBitrate  int    `envconfig:"SRT_MAX_BITRATE" default:"10000000"` // 10 Mbps
		Latency     int    `envconfig:"SRT_LATENCY" default:"120"`         // milliseconds
	}
	
	HLS struct {
		SegmentDuration int    `envconfig:"HLS_SEGMENT_DURATION" default:"6"`
		PlaylistLength  int    `envconfig:"HLS_PLAYLIST_LENGTH" default:"10"`
		StoragePath     string `envconfig:"HLS_STORAGE_PATH" default:"/tmp/hls"`
	}
	
	WebRTC struct {
		STUNServers []string `envconfig:"WEBRTC_STUN_SERVERS" default:"stun:stun.l.google.com:19302"`
		TURNServer  string   `envconfig:"WEBRTC_TURN_SERVER"`
		TURNUser    string   `envconfig:"WEBRTC_TURN_USER"`
		TURNPass    string   `envconfig:"WEBRTC_TURN_PASS"`
	}
	
	// Security
	Security struct {
		JWTSecret      string   `envconfig:"JWT_SECRET" required:"true"`
		AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS" default:"https://app.middleman.com,https://middleman.com"`
		APIKeyHeader   string   `envconfig:"API_KEY_HEADER" default:"X-API-Key"`
	}
	
	// Webhooks
	Webhooks struct {
		MaxRetries       int           `envconfig:"WEBHOOK_MAX_RETRIES" default:"3"`
		Timeout          time.Duration `envconfig:"WEBHOOK_TIMEOUT" default:"30s"`
		Workers          int           `envconfig:"WEBHOOK_WORKERS" default:"10"`
		QueueSize        int           `envconfig:"WEBHOOK_QUEUE_SIZE" default:"1000"`
		HTTPSOnly        bool          `envconfig:"WEBHOOK_HTTPS_ONLY" default:"true"`
		BackoffFactor    float64       `envconfig:"WEBHOOK_BACKOFF_FACTOR" default:"2.0"`
		InitialDelay     time.Duration `envconfig:"WEBHOOK_INITIAL_DELAY" default:"1s"`
		MaxBackoff       time.Duration `envconfig:"WEBHOOK_MAX_BACKOFF" default:"5m"`
		CleanupAge       time.Duration `envconfig:"WEBHOOK_CLEANUP_AGE" default:"720h"` // 30 days
	}
	
	// Database
	Database struct {
		URL             string        `envconfig:"DATABASE_URL" required:"true"`
		MaxConnections  int           `envconfig:"DATABASE_MAX_CONNECTIONS" default:"100"`
		MaxIdleConns    int           `envconfig:"DATABASE_MAX_IDLE_CONNS" default:"10"`
		ConnMaxLifetime time.Duration `envconfig:"DATABASE_CONN_MAX_LIFETIME" default:"1h"`
	}
	
	// CDN Configuration
	CDN struct {
		Provider      string `envconfig:"CDN_PROVIDER" default:"cloudflare"`
		APIKey        string `envconfig:"CDN_API_KEY"`
		ZoneID        string `envconfig:"CDN_ZONE_ID"`
		BaseURL       string `envconfig:"CDN_BASE_URL"`
		PurgeEnabled  bool   `envconfig:"CDN_PURGE_ENABLED" default:"true"`
	}
	
	// DRM Configuration
	DRM struct {
		Enabled           bool   `envconfig:"DRM_ENABLED" default:"false"`
		WidevineURL       string `envconfig:"DRM_WIDEVINE_URL"`
		WidevineKey       string `envconfig:"DRM_WIDEVINE_KEY"`
		FairPlayURL       string `envconfig:"DRM_FAIRPLAY_URL"`
		FairPlayCert      string `envconfig:"DRM_FAIRPLAY_CERT"`
		PlayReadyURL      string `envconfig:"DRM_PLAYREADY_URL"`
		PlayReadyKey      string `envconfig:"DRM_PLAYREADY_KEY"`
	}
	
	// Limits
	Limits struct {
		MaxConcurrentViewers int           `envconfig:"MAX_CONCURRENT_VIEWERS" default:"10000"`
		MaxStreamDuration    time.Duration `envconfig:"MAX_STREAM_DURATION" default:"8h"`
		MaxBitrate           int           `envconfig:"MAX_BITRATE" default:"8000000"` // 8 Mbps
		MaxResolution        string        `envconfig:"MAX_RESOLUTION" default:"1920x1080"`
	}
	
	// Monitoring
	Monitoring struct {
		MetricsEnabled     bool   `envconfig:"METRICS_ENABLED" default:"true"`
		MetricsPort        string `envconfig:"METRICS_PORT" default:"9100"`
		TracingEnabled     bool   `envconfig:"TRACING_ENABLED" default:"true"`
		TracingEndpoint    string `envconfig:"TRACING_ENDPOINT"`
		LogLevel           string `envconfig:"LOG_LEVEL" default:"info"`
		LogFormat          string `envconfig:"LOG_FORMAT" default:"json"`
	}
	
	// Storage
	Storage struct {
		Type          string `envconfig:"STORAGE_TYPE" default:"local"` // local, s3, gcs
		LocalPath     string `envconfig:"STORAGE_LOCAL_PATH" default:"/var/lib/streams"`
		S3Bucket      string `envconfig:"STORAGE_S3_BUCKET"`
		S3Region      string `envconfig:"STORAGE_S3_REGION"`
		S3Endpoint    string `envconfig:"STORAGE_S3_ENDPOINT"`
		GCSBucket     string `envconfig:"STORAGE_GCS_BUCKET"`
		GCSProject    string `envconfig:"STORAGE_GCS_PROJECT"`
	}
	
	// Rate Limiting
	RateLimit struct {
		Enabled        bool   `envconfig:"RATE_LIMIT_ENABLED" default:"true"`
		RequestsPerMin int    `envconfig:"RATE_LIMIT_REQUESTS_PER_MIN" default:"60"`
		BurstSize      int    `envconfig:"RATE_LIMIT_BURST_SIZE" default:"10"`
		RedisURL       string `envconfig:"RATE_LIMIT_REDIS_URL"`
	}
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("STREAMS", &cfg)
	if err != nil {
		return nil, err
	}
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	
	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Add validation logic here
	// For example, check required fields, validate URLs, etc.
	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}