package main

import (
	"database/sql"
	"fmt"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"middleman/managers"
	"middleman/managers/migrations"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {

	if err := run(); err != nil {
		fmt.Printf("managers exitted abnormally: %s\n", err)
		os.Exit(1)
	}

}

func run() (err error) {

	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {

		return fmt.Errorf("failed to initialize app config: %w", err)
	}
	var managerCfg config.ManagersConfig
	managerCfg, err = config.InitManagersConfig()
	if err != nil {

		return
	}

	publicMethods := []string{}

	if err != nil {
		return err
	}
	s, err := system.NewManagersSystem(cfg, publicMethods, managerCfg)
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
	if err = managers.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	fmt.Println("started managers service")
	defer fmt.Println("stopped managers service")

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
