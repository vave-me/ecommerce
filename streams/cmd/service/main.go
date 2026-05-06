package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"middleman/streams"
	"middleman/streams/migrations"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("streams exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

// AI ADD PRODUCT BASE ON THIS IMAGE
func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	publicMethods := []string{
		"/streamspb.StreamsService/GetStreams",
		"/streamspb.StreamsService/GetStreamsWithFilters",
		"/streamspb.StreamsService/GetStreamsByCategory",
		"/streamspb.StreamsService/GetStreamsByCategorySlug",
		"/streamspb.StreamsService/GetStream",
		"/streamspb.StreamsService/GetVariants",
		"/streamspb.StreamsService/GetPublicCatalog",
		// Add other public methods if any
	}
	if err != nil {
		return err
	}
	s, err := system.NewSystem(cfg, publicMethods)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			return
		}
	}(s.DB())
	if err = s.MigrateDB(migrations.FS); err != nil {
		return err
	}
	s.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))
	// call the module composition root
	if err = streams.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	fmt.Println("started streams service")
	defer fmt.Println("stopped streams service")

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

	return s.Waiter().Wait()
}
