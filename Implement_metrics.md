# Integrating our Custom Metrics Package into Hkoffee Go Application

This document provides step-by-step instructions for integrating a custom Prometheus-based metrics package into a typical Go application. The package enables you to collect HTTP metrics (including response times, request counts, errors, and status codes), system metrics, database metrics, and more—all using a custom registry that excludes default Go collectors.

---

## Project Structure

A common Go project structure might look like:

```
myapp/
├── cmd/
│   └── myapp/
│       └── main.go
├── internal/
│   ├── db/
│   │   └── db.go
│   ├── handlers/
│   │   └── handlers.go
│   └── metrics/
│       └── metrics.go
├── go.mod
└── go.sum
```

---

## 1. Create the Metrics Package

Place your custom metrics package in `internal/metrics/metrics.go`. This package should:

Have the latest Metrics package we will finalize.

---

## 2. Create Your HTTP Handlers

Place your HTTP handlers in `internal/handlers/handlers.go`. For example:

```go
package handlers

import (
	"encoding/json"
	"net/http"
)

// HomeHandler responds with a welcome message.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Welcome to MyApp"}
	json.NewEncoder(w).Encode(response)
}

// HealthHandler returns a simple health status.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
```

---

## 3. Database Integration (Optional)

If you have a database, place your connection and query logic in `internal/db/db.go` and instrument it with your metrics package. For example:

```go
package db

import (
	"context"
	"log"

	"myapp/internal/metrics"
	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func Connect(databaseURL string) {
	var err error
	Conn, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	log.Println("Connected to the database")
	metrics.UpdateDBConnections(1) // Update this metric as needed.
}

func Close() {
	if Conn != nil {
		Conn.Close(context.Background())
		log.Println("Disconnected from the database")
		metrics.UpdateDBConnections(0)
	}
}
```

---

## 4. Integrate in main.go

Place your main function in `cmd/myapp/main.go`:

```go
package main

import (
	"log"
	"net/http"

	"myapp/internal/db"
	"myapp/internal/handlers"
	"myapp/internal/metrics"
)

func main() {
	// Register custom metrics.
	metrics.RegisterMetrics()

	// (Optional) Start any background routines for uptime, system metrics, etc.
	go metrics.StartUptime()
	go metrics.UpdateSystemMetrics()

	// Connect to the database.
	db.Connect("postgres://user:password@localhost:5432/mydb")
	defer db.Close()

	// Instrument your HTTP handlers.
	http.Handle("/", metrics.InstrumentHandler("/", http.HandlerFunc(handlers.HomeHandler)))
	http.Handle("/auth/{email}", metrics.InstrumentHandler("/auth/{email}", http.HandlerFunc(handlers.EmailFinder)))
	http.Handle("/health", metrics.InstrumentHandler("/health", http.HandlerFunc(handlers.HealthHandler)))
	// Add other routes as needed.

	// Expose the /metrics endpoint.
	http.Handle("/metrics", metrics.MetricsHandler())

	log.Printf("Server is running on port %s...", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
```

---

## 5. Running the Application and Access Metrics

- Build and run your application:
  ```sh
  go build ./...
  ./myapp
  ```

- Access your application endpoints:
  - `http://localhost:8080/` – Home page.
  - `http://localhost:8080/health` – Health check.

- Check the metrics by visiting:
  ```sh
  curl http://localhost:8080/metrics
  ```
  You should see your custom metrics, including HTTP response times, counts, errors, and status codes.

---

## Conclusion

This guide demonstrates how to integrate a custom metrics package into a typical Go application. By following these steps, you can monitor HTTP performance, system resources, and database operations using Prometheus. Feel free to extend and adjust the package as needed for our specific requirements.