package main

import (
	"log"
	"net/http"

	"federal-funds-rate-metrics-ByYear/config"
	"federal-funds-rate-metrics-ByYear/db"
	"federal-funds-rate-metrics-ByYear/handle"
	"federal-funds-rate-metrics-ByYear/metrics"
)

func main() {
	port := ":8085"
	// Load the configuration from the .env file
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Pass the configuration to packages that need it
	handle.InitConfig(appConfig)

	// Connect to the database
	db.Connect(appConfig.DatabaseURL)
	defer db.Close()

	// Register metrics with Prometheus
	metrics.RegisterMetrics()

	// Start background routines to update metrics.
	metrics.StartUptime()
	metrics.UpdateSystemMetrics()

	// Instrument and register HTTP handlers with static route patterns.
	// These static patterns ensure dynamic parts (e.g., email or id) are not included in the metric labels.
	http.Handle("/", metrics.InstrumentHandler("/", http.HandlerFunc(handle.FederalFundsHandlerInsight)))
	http.Handle("/auth/{email}", metrics.InstrumentHandler("/auth/{email}", http.HandlerFunc(handle.UserInfo)))
	http.Handle("/create", metrics.InstrumentHandler("/create", http.HandlerFunc(handle.CreateUser)))
	http.Handle("/health", metrics.InstrumentHandler("/health", http.HandlerFunc(handle.HealthCheck)))

	// Expose the /metrics endpoint for Prometheus to scrape real-time metrics.
	http.Handle("/metrics", metrics.MetricsHandler())

	log.Printf("Server is running on port %s...", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
