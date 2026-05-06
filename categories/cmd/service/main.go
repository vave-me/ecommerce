package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"middleman/categories"
	"middleman/categories/migrations"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("categories exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	publicMethods := []string{
		"/categoriespb.CategoriesService/GetCategories",
		"/categoriespb.CategoriesService/GetMainCategories",
		"/categoriespb.CategoriesService/GetCategoryBySlug",
		"/categoriespb.CategoriesService/GetAllMainCategories",
		"/categoriespb.CategoriesService/GetSubCategories",
		"/categoriespb.CategoriesService/GetCategoriesBySlug",
		"/categoriespb.CategoriesService/GetCategory",
		"/categoriespb.CategoriesService/GetFilters",
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
		fmt.Printf("migration failed: %s\n", err)
		return err
	}
	s.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))
	// call the module composition root
	if err = categories.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	fmt.Println("started categories service")
	defer fmt.Println("stopped categories service")

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
