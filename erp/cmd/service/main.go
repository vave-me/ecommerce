package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"middleman/erp"
	"middleman/erp/migrations"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("erp service exited abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	// Initialize logger for startup
	logger := zerolog.New(os.Stdout).With().
		Str("service", "erp").
		Timestamp().
		Logger()

	logger.Info().Msg("starting ERP service initialization")

	var cfg config.AppConfig
	logger.Debug().Msg("loading application configuration")
	cfg, err = config.InitConfig()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load configuration")
		return err
	}
	logger.Info().
		Str("environment", cfg.Environment).
		Str("log_level", cfg.LogLevel).
		Str("rpc_address", cfg.Rpc.Address()).
		Str("web_address", cfg.Web.Address()).
		Msg("configuration loaded successfully")

	publicMethods := []string{
		"/api/erp/webhook",
		"/erp/webhook",
	}
	logger.Debug().Strs("methods", publicMethods).Msg("configured public methods")

	logger.Info().Msg("initializing system components")
	s, err := system.NewSystem(cfg, publicMethods)
	if err != nil {
		logger.Error().Err(err).Msg("failed to initialize system")
		return err
	}

	// Use system logger from here
	sysLogger := s.Logger()
	sysLogger.Info().Msg("system initialized successfully")

	defer func(db *sql.DB) {
		sysLogger.Debug().Msg("closing database connection")
		if err = db.Close(); err != nil {
			sysLogger.Error().Err(err).Msg("error closing database connection")
			return
		}
		sysLogger.Debug().Msg("database connection closed")
	}(s.DB())
	sysLogger.Info().Msg("running database migrations")
	startTime := time.Now()
	if err = s.MigrateDB(migrations.FS); err != nil {
		sysLogger.Error().Err(err).Msg("failed to run database migrations")
		return err
	}
	sysLogger.Info().
		Dur("duration", time.Since(startTime)).
		Msg("database migrations completed successfully")
	sysLogger.Debug().Msg("mounting web UI file server")
	s.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))

	// call the module composition root
	sysLogger.Info().Msg("initializing ERP module")
	startTime = time.Now()
	if err = erp.Root(s.Waiter().Context(), s); err != nil {
		sysLogger.Error().Err(err).Msg("failed to initialize ERP module")
		return err
	}
	sysLogger.Info().
		Dur("duration", time.Since(startTime)).
		Msg("ERP module initialized successfully")

	sysLogger.Info().
		Str("rpc_address", cfg.Rpc.Address()).
		Str("web_address", cfg.Web.Address()).
		Str("nats_url", cfg.Nats.URL).
		Msg("ERP service started successfully")

	defer func() {
		sysLogger.Info().Msg("stopping ERP service")
	}()

	sysLogger.Debug().Msg("registering service waiters")
	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
	)
	sysLogger.Info().Msg("all service components registered")

	// Uncomment to enable memory usage monitoring
	// go func() {
	// 	for {
	// 		var mem runtime.MemStats
	// 		runtime.ReadMemStats(&mem)
	// 		sysLogger.Debug().
	// 			Uint64("alloc_kb", mem.Alloc/1024).
	// 			Uint64("total_alloc_kb", mem.TotalAlloc/1024).
	// 			Uint64("sys_kb", mem.Sys/1024).
	// 			Uint32("num_gc", mem.NumGC).
	// 			Msg("memory stats")
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	sysLogger.Info().Msg("waiting for service shutdown signal")
	err = s.Waiter().Wait()
	if err != nil {
		sysLogger.Error().Err(err).Msg("service terminated with error")
	} else {
		sysLogger.Info().Msg("service shutdown completed successfully")
	}
	return err
}
