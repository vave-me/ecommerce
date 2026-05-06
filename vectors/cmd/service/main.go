package main

import (
	"fmt"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	vector "middleman/vectors"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("vector exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}
	var vectorCfg config.VectorsConfig
	vectorCfg, err = config.InitVectorsConfig()
	if err != nil {
		return err
	}

	publicMethods := []string{}
	s, err := system.NewVectorsSystem(cfg, publicMethods, vectorCfg)
	if err != nil {
		return err
	}

	s.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))
	// call the module composition root
	if err = vector.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	fmt.Println("started vector service")
	defer fmt.Println("stopped vector service")

	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
	)

	return s.Waiter().Wait()
}
