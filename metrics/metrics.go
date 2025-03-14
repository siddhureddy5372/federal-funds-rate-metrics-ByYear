package metrics

import (
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Create a custom registry to avoid default collectors.
var customRegistry = prometheus.NewRegistry()

// -----------------------
// Metric Definitions
// -----------------------

// HTTPDuration is a histogram that tracks the response time of HTTP requests.
var (
	HTTPDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_duration_seconds",
		Help:    "Histogram of response time for HTTP requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"path", "method"})
)

// HTTPRequestCount counts the total number of HTTP requests.
var (
	HTTPRequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"path", "method"})

	// HTTPErrorCount counts the total number of HTTP errors.
	HTTPErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "Total number of HTTP errors",
	}, []string{"path", "method"})
)

// HTTPStatusCounter counts the total number of responses for each HTTP status code.
var (
	HTTPStatusCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_status_total",
		Help: "Total number of HTTP responses by status code",
	}, []string{"path", "method", "code"})
)

// Uptime and downtime gauges.
var (
	Uptime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "application_uptime_seconds",
		Help: "Application uptime in seconds",
	})

	Downtime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "application_downtime_seconds",
		Help: "Application downtime in seconds",
	})
)

// AvailabilityRate gauge computed as uptime / (uptime+downtime)*100.
var AvailabilityRate = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "application_availability_rate",
	Help: "Application availability rate in percentage",
})

// CPU and Memory usage gauges.
var (
	CPUUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage_percent",
		Help: "CPU usage percentage",
	})
	MemoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage_bytes",
		Help: "Memory usage in bytes",
	})
)

// Queue length gauge.
var QueueLength = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "queue_length",
	Help: "Length of the processing queue",
})

// Database metrics.
var (
	DBQueryDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "db_query_duration_seconds",
		Help:    "Duration of database queries",
		Buckets: prometheus.DefBuckets,
	})
	DBOpenConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_open_connections",
		Help: "Number of open database connections",
	})
)

// Throughput counter for processed items.
var ThroughputCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "throughput_total",
	Help: "Total number of processed items (throughput)",
})

// -----------------------
// Initialization & Update Functions
// -----------------------

// RegisterMetrics registers all the defined metrics with the custom registry.
func RegisterMetrics() {
	customRegistry.MustRegister(HTTPDuration)
	customRegistry.MustRegister(HTTPRequestCount)
	customRegistry.MustRegister(HTTPErrorCount)
	customRegistry.MustRegister(HTTPStatusCounter)
	customRegistry.MustRegister(Uptime)
	customRegistry.MustRegister(Downtime)
	customRegistry.MustRegister(AvailabilityRate)
	customRegistry.MustRegister(CPUUsage)
	customRegistry.MustRegister(MemoryUsage)
	customRegistry.MustRegister(QueueLength)
	customRegistry.MustRegister(DBQueryDuration)
	customRegistry.MustRegister(DBOpenConnections)
	customRegistry.MustRegister(ThroughputCounter)
}

// MetricsHandler returns an HTTP handler for exposing metrics from the custom registry.
func MetricsHandler() http.Handler {
	return promhttp.HandlerFor(customRegistry, promhttp.HandlerOpts{})
}

// StartUptime begins a background routine to update the uptime gauge.
func StartUptime() {
	startTime := time.Now()
	go func() {
		for {
			uptime := time.Since(startTime).Seconds()
			Uptime.Set(uptime)
			// For demonstration, assume downtime is 0.
			Downtime.Set(0)
			availability := 100.0
			if uptime > 0 {
				availability = (uptime / uptime) * 100.0
			}
			AvailabilityRate.Set(availability)
			time.Sleep(1 * time.Second)
		}
	}()
}

// UpdateSystemMetrics updates CPU and memory usage periodically.
func UpdateSystemMetrics() {
	go func() {
		for {
			// Replace dummy CPU usage with actual metrics if available.
			CPUUsage.Set(10.0)
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			MemoryUsage.Set(float64(m.Alloc))
			time.Sleep(10 * time.Second)
		}
	}()
}

// RecordHTTPMetrics is a helper function to record HTTP metrics.
// It should be called with the request's route pattern, HTTP method, the start time,
// and a flag indicating whether the request resulted in an error.
func RecordHTTPMetrics(path, method string, start time.Time, isError bool) {
	duration := time.Since(start).Seconds()
	HTTPDuration.WithLabelValues(path, method).Observe(duration)
	HTTPRequestCount.WithLabelValues(path, method).Inc()
	if isError {
		HTTPErrorCount.WithLabelValues(path, method).Inc()
	}
}

// RecordDBQuery records a database query duration.
func RecordDBQuery(duration time.Duration) {
	DBQueryDuration.Observe(duration.Seconds())
}

// UpdateDBConnections sets the current number of open database connections.
func UpdateDBConnections(conns int) {
	DBOpenConnections.Set(float64(conns))
}

// -----------------------
// HTTP Metrics Middleware for net/http
// -----------------------

// responseWriter is a wrapper around http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader intercepts the WriteHeader call to capture the HTTP status code.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// HTTPMetricsMiddlewareWithPattern returns a middleware that wraps HTTP requests
// and records metrics using a user-supplied route pattern.
func HTTPMetricsMiddlewareWithPattern(pattern string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)
			// Determine if this is an error.
			isError := rw.statusCode >= http.StatusBadRequest
			// Record basic HTTP metrics.
			RecordHTTPMetrics(pattern, r.Method, start, isError)
			// Record the HTTP status code in a separate counter.
			HTTPStatusCounter.WithLabelValues(pattern, r.Method, strconv.Itoa(rw.statusCode)).Inc()
		})
	}
}

// InstrumentHandler wraps an http.Handler with HTTP metrics instrumentation using the provided route pattern.
func InstrumentHandler(pattern string, handler http.Handler) http.Handler {
	return HTTPMetricsMiddlewareWithPattern(pattern)(handler)
}
