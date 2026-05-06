package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"middleman/merchant"
	"middleman/merchant/migrations"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting merchant service...")
	if err := run(); err != nil {
		fmt.Printf("merchant exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig

	fmt.Println("Initializing application config...")
	cfg, err = config.InitConfig()
	if err != nil {
		fmt.Printf("Failed to initialize app config: %v\n", err)
		return err
	}

	var merchantCfg config.MerchantConfig
	fmt.Println("Initializing merchant config...")
	merchantCfg, err = config.InitMerchantConfig()
	if err != nil {
		fmt.Printf("Failed to initialize merchant config: %v\n", err)
		return err
	}

	// Log merchant configuration (without sensitive data)
	fmt.Printf("Merchant config loaded - MerchantID: %d\n", merchantCfg.MerchantId)
	if merchantCfg.ServiceAccountJSONPath == "" || merchantCfg.ServiceAccountJSONPath == "/path/to/some.json" {
		fmt.Println("WARNING: No valid service account path configured - merchant service will run in disabled mode")
	} else {
		fmt.Printf("Service account path configured: %s\n", merchantCfg.ServiceAccountJSONPath)
	}

	publicMethods := []string{}
	fmt.Println("Creating merchant system...")
	s, err := system.NewMerchantSystem(cfg, publicMethods, merchantCfg)
	if err != nil {
		fmt.Printf("Failed to create merchant system: %v\n", err)
		return err
	}
	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			return
		}
	}(s.DB())
	fmt.Println("Running database migrations...")
	if err = s.MigrateDB(migrations.FS); err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		return err
	}
	fmt.Println("Database migrations completed successfully")
	s.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))
	// call the module composition root
	fmt.Println("Initializing merchant module...")
	if err = merchant.Root(s.Waiter().Context(), s); err != nil {
		fmt.Printf("Failed to initialize merchant module: %v\n", err)
		return err
	}
	fmt.Println("Merchant module initialized successfully")

	fmt.Printf("Merchant service started successfully\n")
	fmt.Printf("- HTTP server listening on: %s\n", cfg.Web.Address())
	fmt.Printf("- gRPC server listening on: %s\n", cfg.Rpc.Address())
	defer fmt.Println("Shutting down merchant service...")

	fmt.Println("Service ready, waiting for shutdown signal...")
	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
	)

	// go func() {
	// 	for {
	// 		var mem runtime.MemStats
	// 		runtime.ReadMemStats(&mem)
	// 		m.logger.Debug().Msgf("Alloc = %v  TotalAlloc = %v  Sys = %v  NumGC = %v", mem.Alloc/1024, mem.TotalAlloc/1024, mem.Sys/1024, mem.NumGC)
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	err = s.Waiter().Wait()
	if err != nil {
		fmt.Printf("Service stopped with error: %v\n", err)
	} else {
		fmt.Println("Service stopped gracefully")
	}
	return err
}
