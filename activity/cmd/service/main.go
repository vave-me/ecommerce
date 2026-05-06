package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
	"middleman/activity"
	"middleman/activity/migrations"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"net/http"
)

func main() {
	// Initialize the logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Run the application and handle errors
	if err := run(logger); err != nil {
		logger.Fatalf("Application exited abnormally: %v", err)
	}
}

func run(logger *logrus.Logger) (err error) {
	// Initialize configuration
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to initialize configuration")
		return err
	}

	var activityCfg config.ActivityConfig
	activityCfg, err = config.InitActivityConfig()

	// Initialize system
	publicMethods := []string{}
	s, err := system.NewActivitySystem(cfg, publicMethods, activityCfg)
	if err != nil {
		logger.WithError(err).Error("Failed to initialize system")
		return err
	}

	// Ensure the database is closed properly
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			logger.WithError(err).Error("Failed to close the database")
		} else {
			logger.Info("Database connection closed successfully")
		}
	}(s.DB())

	if err := s.MigrateDB(migrations.FS); err != nil {
		logger.WithError(err).Error("Database migration failed")
		return err
	}

	// Mount the web UI
	s.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))
	logger.Info("Web UI mounted at '/'")

	// Initialize activity module
	if err := activity.Root(s.Waiter().Context(), s); err != nil {
		logger.WithError(err).Error("Failed to initialize activity module")
		return err
	}

	logger.Info("Activity service started")
	defer logger.Info("Activity service stopped")

	// Add wait functions
	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
	)

	// Optional: Memory stats logging (uncomment if needed)
	/*
		go func() {
			for {
				var mem runtime.MemStats
				runtime.ReadMemStats(&mem)
				logger.Debugf("Alloc = %v KB, TotalAlloc = %v KB, Sys = %v KB, NumGC = %v",
					mem.Alloc/1024, mem.TotalAlloc/1024, mem.Sys/1024, mem.NumGC)
				time.Sleep(10 * time.Second)
			}
		}()
	*/

	// Wait for all services to finish
	if err := s.Waiter().Wait(); err != nil {
		logger.WithError(err).Error("Service encountered an error during wait")
		return err
	}

	return nil
}
