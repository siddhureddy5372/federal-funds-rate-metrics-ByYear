package db

import (
	"context"
	"log"

	"federal-funds-rate-metrics-ByYear/metrics" // import your metrics package

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

// Connect initializes the database connection
func Connect(databaseURL string) {
	var err error
	Conn, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}
	log.Println("Connected to PostgreSQL database!")
	// Update the DB connection metric (for a single connection)
	metrics.UpdateDBConnections(1)
}

// Close terminates the database connection
func Close() {
	if Conn != nil {
		Conn.Close(context.Background())
		log.Println("Disconnected from PostgreSQL database.")
		// Update the DB connection metric to 0 once closed
		metrics.UpdateDBConnections(0)
	}
}
