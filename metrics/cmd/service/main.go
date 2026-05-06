package main

import (
	"database/sql"
	"fmt"
	"middleman/internal/config"
	"middleman/internal/system"
	"middleman/internal/web"
	"middleman/metrics"
	"middleman/metrics/migrations"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("metric exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}
	var metricCfg config.MetricsConfig
	metricCfg, err = config.InitMetricsConfig()
	if err != nil {
		return err
	}

	publicMethods := []string{
		"/metricspb.MetricsService/GetUserMetric",
		"/metricspb.MetricsService/GetItemMetric",
		"/metricspb.MetricsService/GetItemsMetric",
		"/metricspb.MetricsService/GetHighestMetricsByType",
		"/metricspb.MetricsService/GetLowestMetricsByType",
	}
	s, err := system.NewMetricsSystem(cfg, publicMethods, metricCfg)
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
	if err = metric.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	fmt.Println("started metric service")
	defer fmt.Println("stopped metric service")

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
