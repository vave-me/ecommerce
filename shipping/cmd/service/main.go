package main

import (
	"database/sql"
	"log"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"middleman/shipping"
	"middleman/shipping/migrations"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	log.Println("Starting shipping service...")
	if err := run(); err != nil {
		log.Fatalf("Shipping service exited abnormally: %s", err)
		os.Exit(1)
	}
}

// FIND MY SHIPPING
func run() (err error) {
	log.Println("Initializing configuration...")
	var cfg config.AppConfig
	cfg, err = config.InitConfig()

	if err != nil {
		log.Printf("Error initializing configuration: %s", err)
		return err
	}
	log.Println("Configuration initialized successfully.")

	publicMethods := []string{
		"/shippingpb.ShippingService/CreateShipping",
		"/shippingpb.ShippingService/TrackShipping",
	}

	var shippingCfg config.ShippingConfig
	log.Println("Loading shipping configuration...")
	shippingCfg, err = config.InitShippingConfig()
	log.Printf("Loaded shipping configuration: %+v", shippingCfg)

	log.Println("Creating new shipping system...")
	s, err := system.NewShippingSystem(cfg, publicMethods, shippingCfg)
	if err != nil {
		log.Printf("Error creating shipping system: %s", err)
		return err
	}
	log.Println("Shipping system created successfully.")

	defer func(db *sql.DB) {
		log.Println("Closing database connection...")
		if err = db.Close(); err != nil {
			log.Printf("Error closing database connection: %s", err)
			return
		}
		log.Println("Database connection closed successfully.")
	}(s.DB())

	log.Println("Running database migrations...")
	if err = s.MigrateDB(migrations.FS); err != nil {
		log.Printf("Error running database migrations: %s", err)
		return err
	}
	log.Println("Database migrations completed successfully.")

	log.Println("Mounting web server...")
	s.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))
	log.Println("Web server mounted successfully.")

	log.Println("Initializing shipping module...")
	if err = shipping.Root(s.Waiter().Context(), s); err != nil {
		log.Printf("Error initializing shipping module: %s", err)
		return err
	}
	log.Println("Shipping module initialized successfully.")

	log.Println("Starting shipping service...")
	defer log.Println("Shipping service stopped.")

	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
	)

	log.Println("Running service waiter...")
	return s.Waiter().Wait()
}
