package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"middleman/geocoding"
	"middleman/geocoding/migrations"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("geocoding exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}
	var geocodingCfg config.GeocodingConfig
	geocodingCfg, err = config.InitGeocodingConfig()
	if err != nil {
		return err
	}

	publicMethods := []string{
		"/geocodingpb.GeocodingService/SuggestProducts",
		"/geocodingpb.GeocodingService/GeocodingProductsWithFilters",
		"/geocodingpb.GeocodingService/GeocodingProductsWithTerm",
		"/geocodingpb.GeocodingService/SuggestAddress",
		"/geocodingpb.GeocodingService/SuggestCity",
	}
	s, err := system.NewGeocodingSystem(cfg, publicMethods, geocodingCfg)
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
	if err = geocoding.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	fmt.Println("started geocoding service")
	defer fmt.Println("stopped geocoding service")

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
