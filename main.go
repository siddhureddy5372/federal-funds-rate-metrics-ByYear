package main

import (
	"log"
	"net/http"

	"federal-funds-rate-metrics-ByYear/config"
	"federal-funds-rate-metrics-ByYear/db"
	"federal-funds-rate-metrics-ByYear/handle"

	"github.com/gorilla/mux"
)

func main() {
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

	// Initialize router
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", handle.FederalFundsHandlerInsight).Methods("GET")
	r.HandleFunc("/id/{id}", handle.UserInfo).Methods("GET")
	r.HandleFunc("/create", handle.CreateUser).Methods("POST")

	// Start the server
	log.Println("Server is running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", r))
}
