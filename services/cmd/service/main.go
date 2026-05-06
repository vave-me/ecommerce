package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"middleman/services"
	"middleman/services/migrations"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("services exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

// ADD MY SERVICE/BUSISSNES
func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	publicMethods := []string{
		"/servicespb.ServicesService/GetServices",
		"/servicespb.ServicesService/GetPublicCatalog",
		"/servicespb.ServicesService/GetServicesByCategory",
		"/servicespb.ServicesService/GetServicesByCategorySlug",
		"/servicespb.ServicesService/GetServicesWithFilters",
		"/servicespb.ServicesService/GetServicesWithFilter",
		"/servicespb.ServicesService/GetService",
		"/servicespb.ServicesService/GetVariants",
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
	if err = services.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	fmt.Println("started services service")
	defer fmt.Println("stopped services service")

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
