package main

import (
	"database/sql"
	"log"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"middleman/media"
	"middleman/media/migrations"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	log.Println("Starting media service...")
	if err := run(); err != nil {
		log.Fatalf("Media service exited abnormally: %s", err)
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
		"/mediapb.MediaService/GetMedia",
		"/mediapb.MediaService/GetMediaByItem",
		"/mediapb.MediaService/GetAllMediaImages",
		"/mediapb.MediaService/GetAllMediaVideos",
		"/mediapb.MediaService/GetAllItemImages",
		"/mediapb.MediaService/GetAllItemVideos",
		"/mediapb.MediaService/GetAllVideos",
		"/mediapb.MediaService/GetAllImagesByAuthor",
		"/mediapb.MediaService/GetAllVideosByAuthor",
	}

	log.Println("Loading media configuration...")
	mediaCfg := config.LoadMediaConfig()
	log.Printf("Loaded media configuration: %+v", mediaCfg)

	log.Println("Creating new media system...")
	s, err := system.NewMediaSystem(cfg, publicMethods, *mediaCfg)
	if err != nil {
		log.Printf("Error creating media system: %s", err)
		return err
	}
	log.Println("Media system created successfully.")

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

	log.Println("Initializing media module...")
	if err = media.Root(s.Waiter().Context(), s); err != nil {
		log.Printf("Error initializing media module: %s", err)
		return err
	}
	log.Println("Media module initialized successfully.")

	log.Println("Starting media service...")
	defer log.Println("Media service stopped.")

	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
	)

	log.Println("Running service waiter...")
	return s.Waiter().Wait()
}
