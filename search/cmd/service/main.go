package main

import (
	"database/sql"
	"fmt"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"middleman/search"
	"middleman/search/migrations"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("search exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}
	var searchCfg config.SearchConfig
	searchCfg, err = config.InitSearchConfig()
	if err != nil {
		return err
	}

	publicMethods := []string{
		"/searchpb.SearchService/SuggestProducts",
		"/searchpb.SearchService/GetCatalog",
		"/searchpb.SearchService/GetPost",
		"/searchpb.SearchService/GetProduct",
		"/searchpb.SearchService/GetService",
		"/searchpb.SearchService/SearchProductsWithFilters",
		"/searchpb.SearchService/SearchProductsWithFilter",
		"/searchpb.SearchService/SearchProductsWithTerm",
		"/searchpb.SearchService/SearchProductsWithCategorySlug",
		"/searchpb.SearchService/SearchProductsWithCategory",
		"/searchpb.SearchService/SuggestPosts",
		"/searchpb.SearchService/SearchPostsWithFilters",
		"/searchpb.SearchService/SearchPostsWithFilter",
		"/searchpb.SearchService/SearchPostsWithTerm",
		"/searchpb.SearchService/SearchPostsWithCategorySlug",
		"/searchpb.SearchService/SearchPostsWithCategory",
		"/searchpb.SearchService/SuggestServices",
		"/searchpb.SearchService/SearchServicesWithFilter",
		"/searchpb.SearchService/SearchServicesWithFilters",
		"/searchpb.SearchService/SearchServicesWithTerm",
		"/searchpb.SearchService/SearchServicesWithCategorySlug",
		"/searchpb.SearchService/SearchServicesWithCategory",
		"/searchpb.SearchService/UnifiedSearch",
		"/searchpb.SearchService/UnifiedFeed",
	}
	s, err := system.NewSearchSystem(cfg, publicMethods, searchCfg)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		if closeErr := db.Close(); closeErr != nil {
			// Log the error but don't override the main error
			fmt.Printf("error closing database: %v\n", closeErr)
		}
	}(s.DB())
	if err = s.MigrateDB(migrations.FS); err != nil {
		return err
	}
	fmt.Println(s.Redisearch().Info())
	fmt.Println(s.Redisearch().List())
	fmt.Println(s.Redisearch().SynDump("search_index"))
	s.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))
	// call the module composition root
	if err = search.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	fmt.Println("started search service")
	defer fmt.Println("stopped search service")

	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
	)

	return s.Waiter().Wait()
}
