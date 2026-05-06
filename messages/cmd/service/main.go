package main

import (
	"database/sql"
	"fmt"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	msg "middleman/messages"
	"middleman/messages/migrations"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("messages exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}

	var messengerCfg config.MessengerConfig
	messengerCfg, err = config.InitMessengerConfig()
	if err != nil {
		return err
	}

	// Optionally, add a log to confirm the WEBSOCKET URL
	fmt.Printf("WEBSOCKET URL: %s\n", messengerCfg.WEBSOCKET)

	publicMethods := []string{}
	s, err := system.NewMessengerSystem(cfg, publicMethods, messengerCfg)
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
	if err = msg.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	fmt.Println("started messages service")
	defer fmt.Println("stopped messages service")

	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
		s.WaitForWebsocket,
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
