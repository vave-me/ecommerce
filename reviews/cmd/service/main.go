package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"middleman/reviews"
	"middleman/reviews/migrations"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("reviews exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

// ADD A POSITIVE REVIEW
func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}

	publicMethods := []string{
		"/reviewspb.ReviewsService/GetReviews",
		"/reviewspb.ReviewsService/GetReview",
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
	if err = reviews.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	fmt.Println("started reviews service")
	defer fmt.Println("stopped reviews service")

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
