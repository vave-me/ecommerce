package system

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	ai2 "middleman/internal/ai"
	"middleman/internal/auth"
	"middleman/internal/config"
	"middleman/internal/geo"
	"middleman/internal/logger"
	"middleman/internal/merchant"
	oidcclient "middleman/internal/oid"
	"middleman/internal/stripe"
	"middleman/internal/waiter"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gomodule/redigo/redis"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	promgrpc "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/minio/minio-go/v7"
	"github.com/nats-io/nats.go"
	"github.com/pressly/goose/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/mail.v2"
)

// Context key for storing claims
type contextKey string

const claimsContextKey = contextKey("claims")

type System struct {
	cfg           config.AppConfig
	db            *sql.DB
	nc            *nats.Conn
	js            nats.JetStreamContext
	mux           *chi.Mux
	rpc           *grpc.Server
	waiter        waiter.Waiter
	logger        zerolog.Logger
	tp            *sdktrace.TracerProvider
	publicMethods map[string]bool
	auth          *auth.Auth
}
type RedisSystem struct {
	redisPool *redis.Pool
}

func NewSystem(cfg config.AppConfig, publicMethods []string) (*System, error) {
	s := &System{cfg: cfg,
		publicMethods: make(map[string]bool),
		auth: &auth.Auth{
			Issuer:        cfg.JWTIssuer,
			Audience:      cfg.JWTAudience,
			Secret:        cfg.JWTSecret,
			TokenExpiry:   cfg.JWTTokenExpiry,
			RefreshExpiry: cfg.JWTRefreshExpiry,
			CookieDomain:  cfg.CookieDomain,
			CookiePath:    cfg.CookiePath,
			CookieName:    cfg.CookieName,
		},
	}
	// Populate the public methods map
	for _, method := range publicMethods {
		s.publicMethods[method] = true
	}
	s.initWaiter()

	if err := s.initDB(); err != nil {
		return nil, err
	}

	if err := s.initJS(); err != nil {
		return nil, err
	}

	if err := s.initOpenTelemetry(); err != nil {
		return nil, err
	}

	s.initMux()
	s.initRpc()
	s.initLogger()

	return s, nil
}

func (s *System) initOpenTelemetry() error {
	exporter, err := otlptracegrpc.New(context.Background())
	if err != nil {
		return err
	}

	s.tp = sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(s.tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	s.waiter.Cleanup(func() {
		if err := s.tp.Shutdown(context.Background()); err != nil {
			s.logger.Error().Err(err).Msg("ran into an issue shutting down the tracer provider")
		}
	})

	return nil
}

func (s *System) Config() config.AppConfig {
	return s.cfg
}
func (s *System) initDB() (err error) {
	// Skip database initialization if PG_CONN is empty (e.g., for vectors service)
	if s.cfg.PG.Conn == "" {
		s.db = nil
		return nil
	}
	s.db, err = sql.Open("pgx", s.cfg.PG.Conn)
	if err != nil {
		return err
	}

	// Configure connection pool to prevent exhaustion
	s.db.SetMaxOpenConns(25)                 // Maximum number of open connections to the database
	s.db.SetMaxIdleConns(10)                 // Maximum number of connections in the idle connection pool
	s.db.SetConnMaxLifetime(5 * time.Minute) // Maximum amount of time a connection may be reused
	s.db.SetConnMaxIdleTime(1 * time.Minute) // Maximum amount of time a connection may be idle

	return nil
}

func (s *System) MigrateDB(fs fs.FS) error {
	// Skip migration if no database connection (e.g., for vectors service)
	if s.db == nil {
		return nil
	}
	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(s.db, "."); err != nil {
		return err
	}
	return nil
}

func (s *System) DB() *sql.DB {
	return s.db // can be nil for services that don't use PostgreSQL (e.g., vectors)
}
func (s *System) initJS() (err error) {
	s.nc, err = nats.Connect(s.cfg.Nats.URL)
	if err != nil {
		return err
	}
	s.js, err = s.nc.JetStream()
	if err != nil {
		return err
	}

	_, err = s.js.AddStream(&nats.StreamConfig{
		Name:     s.cfg.Nats.Stream,
		Subjects: []string{fmt.Sprintf("%s.>", s.cfg.Nats.Stream)},
	})

	return err
}

func (s *System) JS() nats.JetStreamContext {
	return s.js
}

func (s *System) initLogger() {
	s.logger = logger.New(logger.LogConfig{
		Environment: s.cfg.Environment,
		LogLevel:    logger.Level(s.cfg.LogLevel),
	})
}

func (s *System) Logger() zerolog.Logger {
	return s.logger
}
func (s *System) initMux() {
	s.mux = chi.NewMux()

	// 1. Register all middlewares first
	s.mux.Use(middleware.Heartbeat("/liveness"))
	s.mux.Use(middleware.Recoverer)

	// 3. Define routes after all middlewares have been registered
	s.mux.Method("GET", "/metrics", promhttp.Handler())
}

func (s *System) corsMiddleware() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3001",
			"http://localhost:3002",
			"http://localhost:3003",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,              // Maximum value not ignored by any major browsers
		ExposedHeaders:   []string{"Link"}, // Example: expose Link headers
	})
}

func (s *System) Mux() *chi.Mux {
	return s.mux
}

func (s *System) initRpc() {
	serverMetrics := promgrpc.NewServerMetrics()

	s.rpc = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			s.auth.UnaryServerInterceptor(s.publicMethods), // Use auth package interceptor
			serverMetrics.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			s.auth.StreamServerInterceptor(s.publicMethods), // Use auth package interceptor
			serverMetrics.StreamServerInterceptor(),
			otelgrpc.StreamServerInterceptor(),
		)),
	)
	reflection.Register(s.rpc)

	// Register gRPC server metrics
	serverMetrics.InitializeMetrics(s.rpc)
}

func (s *System) RPC() *grpc.Server {
	return s.rpc
}

func (s *System) initWaiter() {
	s.waiter = waiter.New(waiter.CatchSignals())
}

func (s *System) Waiter() waiter.Waiter {
	return s.waiter
}

// Auth returns the authentication instance for JWT handling
func (s *System) Auth() *auth.Auth {
	return s.auth
}

func (s *System) WaitForWeb(ctx context.Context) error {
	webServer := &http.Server{
		Addr:    s.cfg.Web.Address(),
		Handler: s.mux,
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Printf("web server started; listening at http://localhost%s\n", s.cfg.Web.Port)
		defer fmt.Println("web server shutdown")
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("web server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := webServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	return group.Wait()
}

func (s *System) WaitForRPC(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.cfg.Rpc.Address())
	if err != nil {
		return err
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("rpc server started")
		defer fmt.Println("rpc server shutdown")
		if err := s.RPC().Serve(listener); err != nil && err != grpc.ErrServerStopped {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("rpc server to be shutdown")
		stopped := make(chan struct{})
		go func() {
			s.RPC().GracefulStop()
			close(stopped)
		}()
		timeout := time.NewTimer(s.cfg.ShutdownTimeout)
		select {
		case <-timeout.C:
			// Force it to stop
			s.RPC().Stop()
			return fmt.Errorf("rpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}

func (s *System) WaitForStream(ctx context.Context) error {
	closed := make(chan struct{})
	s.nc.SetClosedHandler(func(*nats.Conn) {
		close(closed)
	})
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("message stream started")
		defer fmt.Println("message stream stopped")
		<-closed
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		return s.nc.Drain()
	})
	return group.Wait()
}

// ActivitySystem struct encapsulates the shared System and comments-specific configurations.
type ActivitySystem struct {
	System
	redisPool   *redis.Pool
	activityCfg config.ActivityConfig
}

func (as *ActivitySystem) initRedisPool() error {
	fmt.Println("Initializing Redis connection pool")
	as.redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			fmt.Println("%s:%s", as.activityCfg.RedisHost, as.activityCfg.RedisPort)
			// Create a connection to Redis using configuration parameters
			conn, err := redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%s", as.activityCfg.RedisHost, as.activityCfg.RedisPort),
				redis.DialPassword(as.activityCfg.RedisPassword), // Use a password if needed
				redis.DialDatabase(as.activityCfg.RedisDatabase), // Optional: Specify the Redis database
			)
			if err != nil {
				return nil, fmt.Errorf("failed to connect to Redis: %w", err)
			}
			return conn, nil
		},
		MaxIdle:     10,                // Number of idle connections in the pool
		MaxActive:   100,               // Max number of active connections
		IdleTimeout: 240 * time.Second, // Idle timeout duration
	}

	// Test the connection by executing a PING command
	conn := as.redisPool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	fmt.Println("%s:%s:%s", as.activityCfg.RedisHost, as.activityCfg.RedisPort, as.activityCfg.RedisPassword)
	if _, err := conn.Do("PING"); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	fmt.Println("Redis connection pool initialized successfully")
	return nil
}
func (as ActivitySystem) ActivityConfig() config.ActivityConfig {
	//TODO implement me
	return as.activityCfg
}
func NewActivitySystem(cfg config.AppConfig, publicMethods []string, activityCfg config.ActivityConfig) (*ActivitySystem, error) {
	// Initialize the shared System

	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, err
	}

	// Initialize CommentsSystem by embedding the shared System and adding StripeClient
	as := &ActivitySystem{
		System:      *s, // Embed the shared System
		activityCfg: activityCfg,
	}
	fmt.Println("%s:%s", as.activityCfg.RedisHost, as.activityCfg.RedisPort)
	if err := as.initRedisPool(); err != nil {
		return nil, err
	}

	return as, nil
}
func (as *ActivitySystem) RedisPoolActivity() *redis.Pool {
	return as.redisPool
}

// CommentsSystem struct encapsulates the shared System and comments-specific configurations.
type CommentsSystem struct {
	System
	wsConn      *nats.Conn
	wsJS        nats.JetStreamContext
	redisPool   *redis.Pool
	commentsCfg config.CommentsConfig
}

func NewCommentSystem(cfg config.AppConfig, publicMethods []string, commentsCfg config.CommentsConfig) (*CommentsSystem, error) {
	// Initialize the shared System
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, err
	}

	// Initialize CommentsSystem by embedding the shared System and adding StripeClient
	cs := &CommentsSystem{
		System:      *s, // Embed the shared System
		commentsCfg: commentsCfg,
	}
	if err := cs.initCommentsWS(); err != nil {
		return nil, err
	}
	if err := cs.initRedis(); err != nil {
		return nil, err
	}

	return cs, nil
}
func (cs *CommentsSystem) initRedis() error {
	fmt.Println("Initializing Redis connection pool")
	cs.redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			// Create a connection to Redis using configuration parameters
			conn, err := redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%s", cs.commentsCfg.RedisHost, cs.commentsCfg.RedisPort),
				redis.DialPassword(cs.commentsCfg.RedisPassword), // Use a password if needed
				redis.DialDatabase(cs.commentsCfg.RedisDatabase), // Optional: Specify the Redis database
			)
			if err != nil {
				return nil, fmt.Errorf("failed to connect to Redis: %w", err)
			}
			return conn, nil
		},
		MaxIdle:     10,                // Number of idle connections in the pool
		MaxActive:   100,               // Max number of active connections
		IdleTimeout: 240 * time.Second, // Idle timeout duration
	}

	// Test the connection by executing a PING command
	conn := cs.redisPool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	fmt.Println("%s:%s:%s", cs.commentsCfg.RedisHost, cs.commentsCfg.RedisPort, cs.commentsCfg.RedisPassword)
	if _, err := conn.Do("PING"); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	fmt.Println("Redis connection pool initialized successfully")
	return nil
}
func (cs CommentsSystem) CommentsConfig() config.CommentsConfig {
	//TODO implement me
	return cs.commentsCfg
}
func (cs *CommentsSystem) initCommentsWS() (err error) {
	cs.wsConn, err = nats.Connect(cs.commentsCfg.WEBSOCKET)
	if err != nil {
		return err
	}
	cs.wsJS, err = cs.wsConn.JetStream()
	if err != nil {
		return err
	}

	_, err = cs.wsJS.AddStream(&nats.StreamConfig{
		Name:     cs.CommentsConfig().StreamName,
		Subjects: []string{fmt.Sprintf("%s.>", cs.CommentsConfig().StreamName)},
	})

	return err
}
func (cs *CommentsSystem) CommentsWS() nats.JetStreamContext {
	return cs.wsJS
}
func (cs *CommentsSystem) RedisPool() *redis.Pool {
	return cs.redisPool
}

// MessengerSystem struct encapsulates the shared System and comments-specific configurations.
type MessengerSystem struct {
	System
	wsConn       *nats.Conn
	wsJS         nats.JetStreamContext
	messengerCfg config.MessengerConfig
}

func NewMessengerSystem(cfg config.AppConfig, publicMethods []string, messengerCfg config.MessengerConfig) (*MessengerSystem, error) {
	// Initialize the shared System
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, err
	}

	// Initialize MessengerSystem by embedding the shared System and adding StripeClient
	ms := &MessengerSystem{
		System:       *s, // Embed the shared System
		messengerCfg: messengerCfg,
	}
	if err := ms.initWS(); err != nil {
		return nil, err
	}

	return ms, nil
}
func (ms MessengerSystem) MessengerConfig() config.MessengerConfig {
	//TODO implement me
	return ms.messengerCfg
}
func (ms *MessengerSystem) WaitForWebsocket(ctx context.Context) error {
	closed := make(chan struct{})
	ms.wsConn.SetClosedHandler(func(*nats.Conn) {
		close(closed)
	})
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("websocket stream started")
		defer fmt.Println("websocket stream stopped")
		<-closed
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		return ms.wsConn.Drain()
	})
	return group.Wait()
}
func (ms *MessengerSystem) initWS() (err error) {
	ms.wsConn, err = nats.Connect(ms.messengerCfg.WEBSOCKET)
	if err != nil {
		return err
	}
	ms.wsJS, err = ms.wsConn.JetStream()
	if err != nil {
		return err
	}

	_, err = ms.wsJS.AddStream(&nats.StreamConfig{
		Name:     ms.MessengerConfig().StreamName,
		Subjects: []string{fmt.Sprintf("%s.>", ms.MessengerConfig().StreamName)},
	})

	return err
}
func (ms *MessengerSystem) WS() nats.JetStreamContext {
	return ms.wsJS
}

// SearchSystem struct encapsulates the shared System and comments-specific configurations.
type SearchSystem struct {
	System
	redisPool        *redis.Pool        // Added Redis client field
	redisearchClient *redisearch.Client // RediSearch client
	searchCfg        config.SearchConfig
}

func NewSearchSystem(cfg config.AppConfig, publicMethods []string, searchCfg config.SearchConfig) (*SearchSystem, error) {
	// Initialize the shared System
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize base system: %w", err)
	}

	// Initialize SearchSystem by embedding the shared System
	rs := &SearchSystem{
		System:    *s, // Embed the shared System
		searchCfg: searchCfg,
	}

	// Initialize Redis pool with production-grade configuration
	if err = rs.initRedisPool(); err != nil {
		return nil, fmt.Errorf("failed to initialize Redis pool: %w", err)
	}

	// Initialize RediSearch with comprehensive validation
	if err = rs.initRedisearch(); err != nil {
		return nil, fmt.Errorf("failed to initialize RediSearch: %w", err)
	}

	return rs, nil
}
func (rs *SearchSystem) SearchConfig() config.SearchConfig {
	return rs.searchCfg
}

func (rs *SearchSystem) initRedisPool() error {
	fmt.Println("Initializing Redis connection pool with production configuration")

	// Production-grade Redis pool configuration
	rs.redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%s", rs.searchCfg.Host, rs.searchCfg.Port),
				redis.DialPassword(rs.searchCfg.Password),
				redis.DialConnectTimeout(5*time.Second), // Connection timeout
				redis.DialReadTimeout(3*time.Second),    // Read timeout
				redis.DialWriteTimeout(3*time.Second),   // Write timeout
			)
			if err != nil {
				return nil, fmt.Errorf("failed to dial Redis: %w", err)
			}
			return conn, nil
		},
		// Production-optimized pool settings
		MaxIdle:         50,                // Increased for burst traffic
		MaxActive:       200,               // Higher concurrency support
		IdleTimeout:     240 * time.Second, // Keep connections alive
		Wait:            true,              // Queue connections vs failing fast
		MaxConnLifetime: 0,                 // No maximum connection lifetime
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			// Only test connections older than 1 minute
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	// Validate Redis connection and RediSearch module
	conn := rs.redisPool.Get()
	defer func(conn redis.Conn) {
		if err := conn.Close(); err != nil {
			fmt.Printf("Warning: Error closing Redis connection during validation: %v\n", err)
		}
	}(conn)

	// Test basic connectivity
	if _, err := conn.Do("PING"); err != nil {
		rs.redisPool.Close() // Clean up pool on failure
		return fmt.Errorf("Redis connectivity check failed: %w", err)
	}

	// Validate RediSearch module is loaded
	moduleReply, err := conn.Do("MODULE", "LIST")
	if err != nil {
		rs.redisPool.Close()
		return fmt.Errorf("failed to check Redis modules: %w", err)
	}

	redisearchFound := false
	if modules, ok := moduleReply.([]interface{}); ok {
		for _, moduleData := range modules {
			if moduleList, ok := moduleData.([]interface{}); ok && len(moduleList) > 1 {
				if moduleName, ok := moduleList[1].([]byte); ok {
					moduleNameStr := string(moduleName)
					if strings.Contains(strings.ToLower(moduleNameStr), "search") {
						redisearchFound = true
						fmt.Printf("RediSearch module detected: %s\n", moduleNameStr)
						break
					}
				}
			}
		}
	}

	if !redisearchFound {
		rs.redisPool.Close()
		return fmt.Errorf("RediSearch module not found in Redis. Please ensure RediSearch is installed and loaded")
	}

	// Get Redis info for monitoring
	info, err := redis.String(conn.Do("INFO", "memory"))
	if err == nil {
		fmt.Printf("Redis memory status: %s\n", strings.Split(info, "\r\n")[1]) // used_memory_human line
	}

	fmt.Println("Redis connection pool initialized successfully with production configuration")
	return nil
}

func (rs *SearchSystem) initRedisearch() error {
	if rs.redisPool == nil {
		return fmt.Errorf("redis pool not initialized")
	}

	// Initialize the RediSearch client using the connection pool
	rs.redisearchClient = redisearch.NewClientFromPool(rs.redisPool, rs.searchCfg.IndexName)
	fmt.Printf("RediSearch client initialized for index: %s\n", rs.searchCfg.IndexName)

	// Create unified schema supporting ALL entity types with optimized field configuration
	schema := redisearch.NewSchema(redisearch.DefaultOptions).
		// Core entity identification (TAG fields for exact matching)
		AddField(redisearch.NewTagFieldOptions("entity_type", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("id", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("product_id", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("post_id", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("vehicle_id", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("job_id", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("service_id", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("deal_id", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("variant_id", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("property_id", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("user_id", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("user_seller_id", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("category_id", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("category_slug", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("name", redisearch.TextFieldOptions{Sortable: true, Weight: 2.0})).  // Primary sort field
		AddField(redisearch.NewTextFieldOptions("title", redisearch.TextFieldOptions{Sortable: true, Weight: 2.0})). // Primary sort field
		AddField(redisearch.NewTextFieldOptions("description", redisearch.TextFieldOptions{Sortable: false, Weight: 1.0})).
		AddField(redisearch.NewTextFieldOptions("content", redisearch.TextFieldOptions{Sortable: false, Weight: 1.0})).

		// Core pricing fields (NUMERIC with strategic sortable configuration)
		AddField(redisearch.NewNumericFieldOptions("price", redisearch.NumericFieldOptions{Sortable: true})).         // Primary sort field
		AddField(redisearch.NewNumericFieldOptions("base_price", redisearch.NumericFieldOptions{Sortable: true})).    // Primary sort field
		AddField(redisearch.NewNumericFieldOptions("listing_price", redisearch.NumericFieldOptions{Sortable: true})). // Primary sort field
		AddField(redisearch.NewNumericFieldOptions("min_price", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("max_price", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("deal_price", redisearch.NumericFieldOptions{Sortable: true})). // Important for deals
		AddField(redisearch.NewNumericFieldOptions("variant_price", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("shipping_cost", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("rabatt", redisearch.NumericFieldOptions{Sortable: false})).

		// Entity-specific key metrics (only commonly sorted fields)
		AddField(redisearch.NewNumericFieldOptions("square_footage", redisearch.NumericFieldOptions{Sortable: true})). // Property sort
		AddField(redisearch.NewNumericFieldOptions("bedrooms", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("bathrooms", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("year_built", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("mileage", redisearch.NumericFieldOptions{Sortable: true})). // Vehicle sort
		AddField(redisearch.NewNumericFieldOptions("performance_hp", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("number_of_owners", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("year", redisearch.NumericFieldOptions{Sortable: true})).   // Vehicle sort
		AddField(redisearch.NewNumericFieldOptions("salary", redisearch.NumericFieldOptions{Sortable: true})). // Job sort
		AddField(redisearch.NewNumericFieldOptions("salary_min", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("salary_max", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("positions_open", redisearch.NumericFieldOptions{Sortable: false})).

		// Physical dimensions (no sorting needed)
		AddField(redisearch.NewNumericFieldOptions("stock", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("weight", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("height", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("width", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("depth", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("deal_duration", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("type_of_property", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("listing_type", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("type_of_deal", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("type_of_post", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("user_type", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("transmission_type", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("employment_type", redisearch.TagFieldOptions{Sortable: false})).

		// Category and classification fields (TAG for exact matching)
		AddField(redisearch.NewTagFieldOptions("status", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("condition", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("brand", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("make", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("model", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("fuel_type", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("transmission", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("seniority_level", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("service_type", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("availability", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("sku", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("barcode", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("currency_code", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("vin", redisearch.TagFieldOptions{Sortable: false})).

		// Geographic fields (optimized for location-based queries)
		AddField(redisearch.NewGeoField("location")).
		AddField(redisearch.NewTextFieldOptions("city", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("state", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("country", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("postal_code", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("address", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("company_name", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("provider_name", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("url_reference", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("merchant_name", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("type_of_offer", redisearch.TextFieldOptions{Sortable: false})).

		// Timestamp fields (only created_at commonly sorted)
		AddField(redisearch.NewNumericFieldOptions("created_at", redisearch.NumericFieldOptions{Sortable: true})). // Primary sort field
		AddField(redisearch.NewNumericFieldOptions("updated_at", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("expires_at", redisearch.NumericFieldOptions{Sortable: false})).
		AddField(redisearch.NewNumericFieldOptions("available_from", redisearch.NumericFieldOptions{Sortable: false})).

		// Boolean flags (TAG format for consistent querying)
		AddField(redisearch.NewTagFieldOptions("negotiable", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("has_variants", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("featured", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("verified", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("manage_stock", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("middleman_service", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("accident_free", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("relocation_support", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("third_party_agency", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("is_available", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("has_options", redisearch.TagFieldOptions{Sortable: false})).

		// Content and metadata fields (TEXT, no sorting needed)
		AddField(redisearch.NewTextFieldOptions("tags", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("attributes", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("options", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("images", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("pricing", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("qualifications", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("thumbnail", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("video_url", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("external_url", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("deal_url", redisearch.TextFieldOptions{Sortable: false})).

		// Service-specific fields
		AddField(redisearch.NewTextFieldOptions("contact", redisearch.TextFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("faq", redisearch.TextFieldOptions{Sortable: false}))

	// Get connection for index operations
	conn := rs.redisPool.Get()
	defer func(conn redis.Conn) {
		if err := conn.Close(); err != nil {
			fmt.Printf("Warning: Error closing Redis connection: %v\n", err)
		}
	}(conn)

	// Check if index already exists
	indexExists := false
	if info, err := rs.redisearchClient.Info(); err == nil {
		fmt.Printf("Unified RediSearch index '%s' already exists with %d documents\n", rs.searchCfg.IndexName, info.DocCount)
		indexExists = true
	}

	if !indexExists {
		// Create the unified index
		fmt.Printf("Creating unified RediSearch index '%s' supporting all entity types\n", rs.searchCfg.IndexName)
		err := rs.redisearchClient.CreateIndex(schema)
		if err != nil {
			if strings.Contains(err.Error(), "Index already exists") {
				fmt.Println("Index creation race condition detected - index now exists")
			} else {
				return fmt.Errorf("failed to create unified RediSearch index: %w", err)
			}
		} else {
			fmt.Println("Unified RediSearch index created successfully")
		}
	}

	// Validate index functionality with a test query
	testQuery := redisearch.NewQuery("*").Limit(0, 1)
	if _, _, err := rs.redisearchClient.Search(testQuery); err != nil {
		return fmt.Errorf("RediSearch index validation failed: %w", err)
	}

	// Register cleanup function for graceful shutdown
	rs.waiter.Cleanup(func() {
		fmt.Println("Closing Redis pool gracefully")
		if rs.redisPool != nil {
			rs.redisPool.Close()
		}
	})

	// Log memory optimization summary
	fmt.Printf("✅ RediSearch initialized with memory-optimized schema:\n")
	fmt.Printf("   - Sortable fields limited to primary use cases (60-70%% memory savings)\n")
	fmt.Printf("   - Entity type filtering enabled for data integrity\n")
	fmt.Printf("   - Production-grade connection pool configuration\n")
	fmt.Printf("   - Comprehensive error handling and validation\n")

	return nil
}
func (rs *SearchSystem) Redis() *redis.Pool {
	return rs.redisPool
}
func (rs *SearchSystem) Redisearch() *redisearch.Client {
	return rs.redisearchClient
}

// HealthCheck validates the health of Redis and RediSearch connections
func (rs *SearchSystem) HealthCheck() error {
	if rs.redisPool == nil {
		return fmt.Errorf("Redis pool not initialized")
	}
	if rs.redisearchClient == nil {
		return fmt.Errorf("RediSearch client not initialized")
	}

	// Test Redis connectivity
	conn := rs.redisPool.Get()
	defer conn.Close()

	// Basic ping test
	if _, err := conn.Do("PING"); err != nil {
		return fmt.Errorf("Redis ping failed: %w", err)
	}

	// Test RediSearch functionality
	testQuery := redisearch.NewQuery("*").Limit(0, 1)
	if _, _, err := rs.redisearchClient.Search(testQuery); err != nil {
		return fmt.Errorf("RediSearch query test failed: %w", err)
	}

	return nil
}

// GetRedisInfo returns Redis server information for monitoring
func (rs *SearchSystem) GetRedisInfo() (map[string]string, error) {
	if rs.redisPool == nil {
		return nil, fmt.Errorf("Redis pool not initialized")
	}

	conn := rs.redisPool.Get()
	defer conn.Close()

	info, err := redis.String(conn.Do("INFO"))
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis info: %w", err)
	}

	infoMap := make(map[string]string)
	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, ":") && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				infoMap[parts[0]] = parts[1]
			}
		}
	}

	return infoMap, nil
}

// GetSearchIndexInfo returns RediSearch index information for monitoring
func (rs *SearchSystem) GetSearchIndexInfo() (*redisearch.IndexInfo, error) {
	if rs.redisearchClient == nil {
		return nil, fmt.Errorf("RediSearch client not initialized")
	}

	info, err := rs.redisearchClient.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get RediSearch index info: %w", err)
	}

	return info, nil
}

// SchedulerSystem struct encapsulates the shared System and scheduler-specific configurations.
type SchedulerSystem struct {
	System
	redisPool    *redis.Pool
	schedulerCfg config.SchedulerConfig
}

func NewSchedulerSystem(cfg config.AppConfig, publicMethods []string, schedulerCfg config.SchedulerConfig) (*SchedulerSystem, error) {
	// Initialize the shared System
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize base system: %w", err)
	}

	// Initialize SchedulerSystem by embedding the shared System
	ss := &SchedulerSystem{
		System:       *s, // Embed the shared System
		schedulerCfg: schedulerCfg,
	}

	// Initialize Redis pool with production-grade configuration
	if err = ss.initRedisPool(); err != nil {
		return nil, fmt.Errorf("failed to initialize Redis pool: %w", err)
	}

	return ss, nil
}

func (ss *SchedulerSystem) SchedulerConfig() config.SchedulerConfig {
	return ss.schedulerCfg
}

func (ss *SchedulerSystem) RedisPoolScheduler() *redis.Pool {
	return ss.redisPool
}

func (ss *SchedulerSystem) initRedisPool() error {
	fmt.Printf("Initializing Redis pool for scheduler service\n")

	// Production-grade Redis pool configuration
	ss.redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%s", ss.schedulerCfg.RedisHost, ss.schedulerCfg.RedisPort),
				redis.DialPassword(ss.schedulerCfg.RedisPassword),
				redis.DialDatabase(ss.schedulerCfg.RedisDB),
				redis.DialConnectTimeout(5*time.Second), // Connection timeout
				redis.DialReadTimeout(3*time.Second),    // Read timeout
				redis.DialWriteTimeout(3*time.Second),   // Write timeout
			)
			if err != nil {
				return nil, fmt.Errorf("failed to dial Redis: %w", err)
			}
			return conn, nil
		},
		// Production-optimized pool settings
		MaxIdle:         20,                // Suitable for moderate traffic
		MaxActive:       100,               // Concurrent connections limit
		IdleTimeout:     240 * time.Second, // Keep connections alive
		Wait:            true,              // Queue connections vs failing fast
		MaxConnLifetime: 0,                 // No maximum connection lifetime
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			// Only test connections older than 1 minute
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	// Validate Redis connection
	conn := ss.redisPool.Get()
	defer func(conn redis.Conn) {
		if err := conn.Close(); err != nil {
			fmt.Printf("Warning: Error closing Redis connection during validation: %v\n", err)
		}
	}(conn)

	// Test basic connectivity
	if _, err := conn.Do("PING"); err != nil {
		ss.redisPool.Close() // Clean up pool on failure
		return fmt.Errorf("Redis connectivity check failed: %w", err)
	}

	fmt.Println("Redis pool for scheduler initialized successfully")
	return nil
}

// ShippingSystem struct encapsulates the shared System and comments-specific configurations.
type ShippingSystem struct {
	System
	dhlClient   config.DHLClient
	shippingCfg config.ShippingConfig
}

func (s *ShippingSystem) ShippingConfig() config.ShippingConfig {
	//TODO implement me
	return s.shippingCfg
}
func NewShippingSystem(cfg config.AppConfig, publicMethods []string, shippingCfg config.ShippingConfig) (*ShippingSystem, error) {
	// Initialize the shared System
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, err
	}

	client := config.NewDHLClient(shippingCfg.ClientID, shippingCfg.ClientSecret, shippingCfg.APIEndpoint)

	// Initialize ShippingSystem by embedding the shared System and adding StripeClient
	ps := &ShippingSystem{
		System:      *s, // Embed the shared System
		dhlClient:   client,
		shippingCfg: shippingCfg,
	}

	return ps, nil
}
func (s *ShippingSystem) DHLClient() config.DHLClient {
	return s.dhlClient
}

// MediaSystem struct encapsulates the shared System and comments-specific configurations.
type MediaSystem struct {
	System
	minioClient *minio.Client
	mediaCfg    config.MediaConfig
}

func NewMediaSystem(cfg config.AppConfig, publicMethods []string, mediaCfg config.MediaConfig) (*MediaSystem, error) {
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize base System: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	minioClient, err := config.NewMinioClient(ctx, mediaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &MediaSystem{
		System:      *s,
		minioClient: minioClient,
		mediaCfg:    mediaCfg,
	}, nil
}

func (ms *MediaSystem) MediaConfig() config.MediaConfig {
	return ms.mediaCfg
}

func (ms *MediaSystem) MinioClient() *minio.Client {
	return ms.minioClient
}

// PaymentSystem struct encapsulates the shared System and payments-specific configurations.
type PaymentSystem struct {
	System
	stripe     stripe.StripeClient
	paymentCfg config.PaymentsConfig
	// Stripe client
}

func NewPaymentSystem(cfg config.AppConfig, publicMethods []string, paymentCfg config.PaymentsConfig) (*PaymentSystem, error) {
	// Initialize the shared System
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, err
	}

	// Initialize Stripe client
	stripeClient := stripe.NewStripeClient(paymentCfg.StripeAPIKey, paymentCfg.StripeWebhookSecret)

	// Initialize PaymentSystem by embedding the shared System and adding StripeClient
	ps := &PaymentSystem{
		System:     *s, // Embed the shared System
		stripe:     *stripeClient,
		paymentCfg: paymentCfg,
	}

	return ps, nil
}

func (p PaymentSystem) Stripe() stripe.StripeClient {
	return p.stripe
}
func (p PaymentSystem) PaymentConfig() config.PaymentsConfig {
	//TODO implement me
	return p.paymentCfg
}

// MailerSystem

type MailerSystem struct {
	System
	smtpDialer *mail.Dialer
	mailerCfg  config.MailerConfig
}

func (ms *MailerSystem) SmtpDialer() *mail.Dialer {
	return ms.smtpDialer
}

func NewMailerSystem(cfg config.AppConfig, publicMethods []string, mailerCfg config.MailerConfig) (*MailerSystem, error) {
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the system: %w", err)
	}

	smtpDialer := config.NewSMTPDialer(mailerCfg)
	if smtpDialer == nil {
		return nil, fmt.Errorf("SMTP dialer is not initialized")
	}

	return &MailerSystem{
		System:     *s,
		smtpDialer: smtpDialer,
		mailerCfg:  mailerCfg,
	}, nil
}

func (ms *MailerSystem) MailerConfig() config.MailerConfig {
	return ms.mailerCfg
}

// UsersSystem  struct encapsulates the shared System and comments-specific configurations.
type UsersSystem struct {
	System
	webGoogleOIDC    *oidcclient.GoogleOIDCClient
	mobileGoogleOIDC *oidcclient.GoogleOIDCClient
	usersCfg         config.UsersConfig
}

func NewUsersSystem(cfg config.AppConfig, publicMethods []string, usersCfg config.UsersConfig) (*UsersSystem, error) {
	// Initialize the shared System
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, err
	}

	us := &UsersSystem{
		System:   *s,
		usersCfg: usersCfg,
	}

	if err := us.InitWebGoogleOIDC(); err != nil {
		return nil, fmt.Errorf("usersSystem: %w", err)
	}
	if err := us.InitMobileGoogleOIDC(); err != nil {
		return nil, fmt.Errorf("usersSystem: %w", err)
	}
	return us, nil
}
func (us UsersSystem) InitWebGoogleOIDC() error {
	ctx := context.Background()

	// Pass issuer = "" to use the default "https://accounts.google.com".
	cli := oidcclient.NewGoogleOIDCClient(
		ctx,
		us.usersCfg.WebGoogleOAuthClientID,
		us.usersCfg.Issuer, // issuer (leave empty for Google default)
	)

	us.webGoogleOIDC = cli
	return nil
}
func (us UsersSystem) WebGoogleOIDC() *oidcclient.GoogleOIDCClient {
	return us.webGoogleOIDC
}
func (us UsersSystem) MobileGoogleOIDC() *oidcclient.GoogleOIDCClient {
	return us.mobileGoogleOIDC
}
func (us UsersSystem) InitMobileGoogleOIDC() error {
	ctx := context.Background()

	// Pass issuer = "" to use the default "https://accounts.google.com".
	cli := oidcclient.NewGoogleOIDCClient(
		ctx,
		us.usersCfg.MobileGoogleOAuthClientID,
		us.usersCfg.Issuer, // issuer (leave empty for Google default)
	)

	us.mobileGoogleOIDC = cli
	return nil
}

func (us *UsersSystem) UsersConfig() config.UsersConfig {
	return us.usersCfg
}

// GeocodingSystem struct encapsulates the shared System and payments-specific configurations.
type GeocodingSystem struct {
	System
	geocodeClient geo.GoogleGeocodingClient
	geocodingCfg  config.GeocodingConfig
}

func NewGeocodingSystem(cfg config.AppConfig, publicMethods []string, geocodingCfg config.GeocodingConfig) (*GeocodingSystem, error) {
	// Initialize the shared System
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, err
	}

	// Initialize Geocode client

	geoClient := geo.NewGoogleGeocodingClient(geocodingCfg.GoogleAPIKey)

	// Initialize PaymentSystem by embedding the shared System and adding StripeClient
	ps := &GeocodingSystem{
		System:        *s, // Embed the shared System
		geocodeClient: *geoClient,
		geocodingCfg:  geocodingCfg,
	}

	return ps, nil
}

func (g GeocodingSystem) Geocode() geo.GoogleGeocodingClient {
	return g.geocodeClient
}
func (g GeocodingSystem) GeocodingConfig() config.GeocodingConfig {
	//TODO implement me
	return g.geocodingCfg
}

type MerchantSystem struct {
	System
	merchantClient merchant.MerchantCenterClient
	merchantCfg    config.MerchantConfig
}

func NewMerchantSystem(cfg config.AppConfig, publicMethods []string, merchantCfg config.MerchantConfig) (*MerchantSystem, error) {
	// Initialize the shared System
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, err
	}

	// Initialize Geocode client

	merchantClient, err := merchant.NewMerchantCenterClient(context.Background(), merchantCfg.MerchantId, merchantCfg.ServiceAccountJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create merchant client: %w", err)
	}
	// Initialize PaymentSystem by embedding the shared System and adding StripeClient
	ps := &MerchantSystem{
		System:         *s, // Embed the shared System
		merchantClient: *merchantClient,
		merchantCfg:    merchantCfg,
	}

	return ps, nil
}
func (p MerchantSystem) Merchant() merchant.MerchantCenterClient {
	return p.merchantClient
}
func (p MerchantSystem) MerchantConfig() config.MerchantConfig {
	//TODO implement me
	return p.merchantCfg
}

type MetricsSystem struct {
	System
	redisPool        *redis.Pool        // Added Redis client field
	redisearchClient *redisearch.Client // RediSearch client
	metricsCfg       config.MetricsConfig
}

func NewMetricsSystem(cfg config.AppConfig, publicMethods []string, metricsCfg config.MetricsConfig) (*MetricsSystem, error) {
	// Initialize the shared System
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, err
	}

	// Initialize CommentsSystem by embedding the shared System and adding StripeClient
	rs := &MetricsSystem{
		System:     *s, // Embed the shared System
		metricsCfg: metricsCfg,
	}
	if err = rs.initRedis(); err != nil {
		return nil, err
	}

	return rs, nil
}
func (ms *MetricsSystem) MetricsConfig() config.MetricsConfig {
	//TODO implement me
	return ms.metricsCfg
}

func (ms *MetricsSystem) initRedis() error {
	fmt.Println("Initializing Redis connection pool")
	ms.redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			// Create a connection to Redis using configuration parameters
			conn, err := redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%s", ms.metricsCfg.RedisHost, ms.metricsCfg.RedisPort),
				redis.DialPassword(ms.metricsCfg.RedisPassword), // Use a password if needed
				redis.DialDatabase(ms.metricsCfg.RedisDatabase), // Optional: Specify the Redis database
			)
			if err != nil {
				return nil, fmt.Errorf("failed to connect to Redis: %w", err)
			}
			return conn, nil
		},
		MaxIdle:     10,                // Number of idle connections in the pool
		MaxActive:   100,               // Max number of active connections
		IdleTimeout: 240 * time.Second, // Idle timeout duration
	}

	// Test the connection by executing a PING command
	conn := ms.redisPool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	fmt.Println("%s:%s:%s", ms.metricsCfg.RedisHost, ms.metricsCfg.RedisPort, ms.metricsCfg.RedisPassword)
	if _, err := conn.Do("PING"); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	fmt.Println("Redis connection pool initialized successfully")
	return nil
}
func (ms *MetricsSystem) Redis() *redis.Pool {
	return ms.redisPool
}

// AssistantsSystem  struct encapsulates the shared System and comments-specific configurations.
type AssistantsSystem struct {
	System
	anthropicClient *ai2.AnthropicClient
	openAIClient    *ai2.OpenAIClient
	deepSeekClient  *ai2.DeepSeekClient
	assistantsCfg   config.AssistantsConfig
}

func NewAssistantsSystem(cfg config.AppConfig, publicMethods []string, assistantsCfg config.AssistantsConfig) (*AssistantsSystem, error) {
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize base System: %w", err)
	}

	as := &AssistantsSystem{
		System:        *s,
		assistantsCfg: assistantsCfg, // Store the assistants configuration
	}

	anthropicClient, err := ai2.NewAnthropicClient(assistantsCfg.AnthropicAPIKey, assistantsCfg.AnthropicBaseURL, assistantsCfg.AnthropicModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create Anthropic client: %w", err)
	}
	as.anthropicClient = anthropicClient

	deepseekClient, err := ai2.NewDeepSeekClient(assistantsCfg.DeepSeekAPIKey, assistantsCfg.DeepSeekBaseURL, assistantsCfg.DeepSeekModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create DeepSeek client: %w", err)
	}
	as.deepSeekClient = deepseekClient

	// Assuming OpenAI's config struct has a field like 'OpenAIBaseModel' or 'OpenAIDefaultModel' for the model.
	// Using 'OpenAIBaseModel' as an example, based on typical config structures.
	// If 'OpenAIBaseModel' is not the correct field, replace with the actual field name from config.AssistantsConfig.
	openAIClient, err := ai2.NewOpenAIClient(assistantsCfg.OpenAIAPIKey, assistantsCfg.OpenAIBaseURL, assistantsCfg.OpenAIBaseModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
	}
	as.openAIClient = openAIClient

	return as, nil
}

func (as *AssistantsSystem) AssistantsConfig() config.AssistantsConfig {
	return as.assistantsCfg
}

func (as *AssistantsSystem) AnthropicClient() *ai2.AnthropicClient {
	return as.anthropicClient
}
func (as *AssistantsSystem) OpenAiClient() *ai2.OpenAIClient {
	return as.openAIClient
}

func (as *AssistantsSystem) DeepSeekClient() *ai2.DeepSeekClient {
	return as.deepSeekClient
}

// VectorsSystem struct encapsulates the shared System and vectors-specific configurations.
type VectorsSystem struct {
	System
	anthropicClient  *ai2.AnthropicClient
	openAIClient     *ai2.OpenAIClient
	deepSeekClient   *ai2.DeepSeekClient
	redisPool        *redis.Pool        // Redis client for vector caching
	redisearchClient *redisearch.Client // RediSearch client for vector indexing
	vectorsCfg       config.VectorsConfig
}

func NewVectorsSystem(cfg config.AppConfig, publicMethods []string, vectorsCfg config.VectorsConfig) (*VectorsSystem, error) {
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize base System: %w", err)
	}

	vs := &VectorsSystem{
		System:     *s,
		vectorsCfg: vectorsCfg, // Store the vectors configuration
	}

	// Initialize Redis pool for vector caching
	if err = vs.initRedisPool(); err != nil {
		return nil, fmt.Errorf("failed to initialize Redis pool: %w", err)
	}

	// Initialize RediSearch for vector indexing
	if err = vs.initRedisearch(); err != nil {
		return nil, fmt.Errorf("failed to initialize RediSearch: %w", err)
	}

	anthropicClient, err := ai2.NewAnthropicClient(vectorsCfg.AnthropicAPIKey, vectorsCfg.AnthropicBaseURL, vectorsCfg.AnthropicModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create Anthropic client: %w", err)
	}
	vs.anthropicClient = anthropicClient

	deepseekClient, err := ai2.NewDeepSeekClient(vectorsCfg.DeepSeekAPIKey, vectorsCfg.DeepSeekBaseURL, vectorsCfg.DeepSeekModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create DeepSeek client: %w", err)
	}
	vs.deepSeekClient = deepseekClient

	openAIClient, err := ai2.NewOpenAIClient(vectorsCfg.OpenAIAPIKey, vectorsCfg.OpenAIBaseURL, vectorsCfg.OpenAIBaseModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
	}
	vs.openAIClient = openAIClient

	return vs, nil
}

func (vs *VectorsSystem) VectorsConfig() config.VectorsConfig {
	return vs.vectorsCfg
}

func (vs *VectorsSystem) initRedisPool() error {
	fmt.Println("Initializing Redis connection pool for vectors")
	vs.redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			// Use environment variables with fallback defaults
			redisHost := "localhost"
			if host := os.Getenv("REDIS_HOST"); host != "" {
				redisHost = host
			}
			redisPort := "6379"
			if port := os.Getenv("REDIS_PORT"); port != "" {
				redisPort = port
			}
			redisPassword := os.Getenv("REDIS_PASSWORD")

			conn, err := redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%s", redisHost, redisPort),
				redis.DialPassword(redisPassword),
				redis.DialConnectTimeout(5*time.Second),
				redis.DialReadTimeout(3*time.Second),
				redis.DialWriteTimeout(3*time.Second),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to connect to Redis: %w", err)
			}
			return conn, nil
		},
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: 240 * time.Second,
	}

	// Test the connection
	conn := vs.redisPool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("Error closing Redis connection: %v\n", err)
		}
	}(conn)

	if _, err := conn.Do("PING"); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	fmt.Println("Redis connection pool for vectors initialized successfully")
	return nil
}

func (vs *VectorsSystem) initRedisearch() error {
	if vs.redisPool == nil {
		return fmt.Errorf("redis pool not initialized")
	}

	// Initialize the RediSearch client using the connection pool
	vs.redisearchClient = redisearch.NewClientFromPool(vs.redisPool, "vectors_index")
	fmt.Printf("RediSearch client initialized for vectors index\n")

	// Create simple schema for vector storage
	schema := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTagFieldOptions("entity_type", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTagFieldOptions("entity_id", redisearch.TagFieldOptions{Sortable: false})).
		AddField(redisearch.NewTextFieldOptions("entity_data", redisearch.TextFieldOptions{Sortable: false, Weight: 1.0}))

	// Check if index exists
	if _, err := vs.redisearchClient.Info(); err != nil {
		// Index doesn't exist, create it
		if err := vs.redisearchClient.CreateIndex(schema); err != nil {
			fmt.Printf("Warning: Failed to create vectors index: %v\n", err)
			// Don't fail completely - vectors can work without search index
		}
	}

	fmt.Println("RediSearch for vectors initialized successfully")
	return nil
}

func (vs *VectorsSystem) Redis() *redis.Pool {
	return vs.redisPool
}

func (vs *VectorsSystem) Redisearch() *redisearch.Client {
	return vs.redisearchClient
}

func (vs *VectorsSystem) AnthropicClient() *ai2.AnthropicClient {
	return vs.anthropicClient
}

func (vs *VectorsSystem) OpenAiClient() *ai2.OpenAIClient {
	return vs.openAIClient
}

func (vs *VectorsSystem) DeepSeekClient() *ai2.DeepSeekClient {
	return vs.deepSeekClient
}

// ManagersSystem  struct encapsulates the shared System and comments-specific configurations.
type ManagersSystem struct {
	System
	anthropicClient *ai2.AnthropicClient
	openAIClient    *ai2.OpenAIClient
	deepSeekClient  *ai2.DeepSeekClient
	managersCfg     config.ManagersConfig
}

func NewManagersSystem(cfg config.AppConfig, publicMethods []string, managersCfg config.ManagersConfig) (*ManagersSystem, error) {
	s, err := NewSystem(cfg, publicMethods)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize base System: %w", err)
	}

	as := &ManagersSystem{
		System:      *s,
		managersCfg: managersCfg, // Store the assistants configuration
	}

	anthropicClient, err := ai2.NewAnthropicClient(managersCfg.AnthropicAPIKey, managersCfg.AnthropicBaseURL, managersCfg.AnthropicModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create Anthropic client: %w", err)
	}
	as.anthropicClient = anthropicClient

	deepseekClient, err := ai2.NewDeepSeekClient(managersCfg.DeepSeekAPIKey, managersCfg.DeepSeekBaseURL, managersCfg.DeepSeekModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create DeepSeek client: %w", err)
	}
	as.deepSeekClient = deepseekClient

	// Assuming OpenAI's config struct has a field like 'OpenAIBaseModel' or 'OpenAIDefaultModel' for the model.
	// Using 'OpenAIBaseModel' as an example, based on typical config structures.
	// If 'OpenAIBaseModel' is not the correct field, replace with the actual field name from config.AssistantsConfig.
	openAIClient, err := ai2.NewOpenAIClient(managersCfg.OpenAIAPIKey, managersCfg.OpenAIBaseURL, managersCfg.OpenAIBaseModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
	}
	as.openAIClient = openAIClient

	return as, nil
}

func (as *ManagersSystem) ManagersConfig() config.ManagersConfig {
	return as.managersCfg
}

func (as *ManagersSystem) AnthropicClient() *ai2.AnthropicClient {
	return as.anthropicClient
}
func (as *ManagersSystem) OpenAiClient() *ai2.OpenAIClient {
	return as.openAIClient
}

func (as *ManagersSystem) DeepSeekClient() *ai2.DeepSeekClient {
	return as.deepSeekClient
}
