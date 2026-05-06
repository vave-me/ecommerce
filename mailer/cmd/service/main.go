package main

import (
	"database/sql"
	"log"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"middleman/mailer"
	"middleman/mailer/migrations"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	log.Println("Starting mailer service...")
	if err := run(); err != nil {
		log.Fatalf("Mailer service exited abnormally: %s", err)
		os.Exit(1)
	}
}

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
		"/mailerpb.MailerService/CreateEmail",
	}

	log.Println("Loading mailer configuration...")
	mailerCfg := config.LoadMailerConfig()
	log.Printf("Loaded mailer configuration: %+v", mailerCfg)

	log.Println("Creating new mailer system...")
	s, err := system.NewMailerSystem(cfg, publicMethods, mailerCfg)
	if err != nil {
		log.Printf("Error creating mailer system: %s", err)
		return err
	}
	log.Println("Mailer system created successfully.")

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

	log.Println("Initializing mailer module...")
	if err = mailer.Root(s.Waiter().Context(), s); err != nil {
		log.Printf("Error initializing mailer module: %s", err)
		return err
	}
	log.Println("Mailer module initialized successfully.")

	log.Println("Starting mailer service...")
	defer log.Println("Mailer service stopped.")

	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
	)

	log.Println("Running service waiter...")
	return s.Waiter().Wait()
}
