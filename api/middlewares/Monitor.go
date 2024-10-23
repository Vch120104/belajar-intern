package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Histogram for request duration in seconds
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler", "method", "status"},
	)

	// Counter for total number of HTTP requests by status code
	requestsByStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by status code.",
		},
		[]string{"handler", "method", "status"},
	)

	// Counters for POST, GET, and PATCH requests
	postCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_post_requests_total",
			Help: "Total number of POST requests.",
		},
		[]string{"status"},
	)

	getCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_get_requests_total",
			Help: "Total number of GET requests.",
		},
		[]string{"status"},
	)

	patchCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_patch_requests_total",
			Help: "Total number of PATCH requests.",
		},
		[]string{"status"},
	)

	// Counter for total number of Prometheus scrapes by HTTP status code
	metricHandlerRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "promhttp_metric_handler_requests_total",
			Help: "Total number of scrapes by HTTP status code.",
		},
		[]string{"code"},
	)

	// Histogram for request and response sizes in bytes
	requestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests in bytes.",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8), // 100B to ~1GB
		},
		[]string{"handler", "method", "status"},
	)

	responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes.",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8), // 100B to ~1GB
		},
		[]string{"handler", "method", "status"},
	)

	// Gauge for the current number of active requests
	activeRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Current number of active requests.",
		},
		[]string{"handler", "method"},
	)

	// Counter for error HTTP requests
	errorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_error_requests_total",
			Help: "Total number of error HTTP requests.",
		},
		[]string{"handler", "method", "status"},
	)

	// Counter for request counts over time
	httpResponseCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_count",
			Help: "Total number of HTTP responses.",
		},
		[]string{"handler", "method", "status"},
	)

	// Histogram for request latency distribution
	httpResponseBucket = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_bucket",
			Help:    "Request latency distribution in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler", "method", "status"},
	)

	// Counter for request counts with intervals
	httpResponseIntervalCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_interval_count",
			Help: "Total number of HTTP responses with interval labels.",
		},
		[]string{"handler", "method", "status", "interval"},
	)

	// SQL Server metrics
	appSQLStatsBucket = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_sql_stats_bucket",
			Help:    "SQL Server query duration bucket.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type"},
	)

	appSQLStatsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_sql_stats_count",
			Help: "SQL Server query count over time.",
		},
		[]string{"query_type"},
	)

	appSQLInUseConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_sql_inUse_connections",
			Help: "Current number of SQL Server in-use connections.",
		},
	)

	appSQLOpenConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_sql_open_connections",
			Help: "Current number of SQL Server open connections.",
		},
	)
)

func init() {
	// Register all metrics with Prometheus
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestsByStatus)
	prometheus.MustRegister(postCounter)
	prometheus.MustRegister(getCounter)
	prometheus.MustRegister(patchCounter)
	prometheus.MustRegister(metricHandlerRequestsTotal)
	prometheus.MustRegister(requestSize)
	prometheus.MustRegister(responseSize)
	prometheus.MustRegister(activeRequests)
	prometheus.MustRegister(errorCounter)
	prometheus.MustRegister(httpResponseCount)
	prometheus.MustRegister(httpResponseBucket)
	prometheus.MustRegister(httpResponseIntervalCount)
	prometheus.MustRegister(appSQLStatsBucket)
	prometheus.MustRegister(appSQLStatsCount)
	prometheus.MustRegister(appSQLInUseConnections)
	prometheus.MustRegister(appSQLOpenConnections)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		activeRequests.WithLabelValues(r.URL.Path, r.Method).Inc()
		defer activeRequests.WithLabelValues(r.URL.Path, r.Method).Dec()

		// Wrap the response writer to capture the status code and response size
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Pass the request through the handler chain
		next.ServeHTTP(rw, r)

		// Calculate duration of the request
		duration := time.Since(start).Seconds()

		// Get the status code of the response
		status := strconv.Itoa(rw.status)

		// Record the duration of the request
		requestDuration.WithLabelValues(r.URL.Path, r.Method, status).Observe(duration)

		// Increment the counter for requests by status code
		requestsByStatus.WithLabelValues(r.URL.Path, r.Method, status).Inc()

		// Increment the method-specific counter
		switch r.Method {
		case http.MethodPost:
			postCounter.WithLabelValues(status).Inc()
		case http.MethodGet:
			getCounter.WithLabelValues(status).Inc()
		case http.MethodPatch:
			patchCounter.WithLabelValues(status).Inc()
		}

		// Increment the promhttp_metric_handler_requests_total counter for specific status codes
		if rw.status == http.StatusNotFound || rw.status == http.StatusCreated {
			metricHandlerRequestsTotal.WithLabelValues(status).Inc()
		}

		// Record the response size
		responseSize.WithLabelValues(r.URL.Path, r.Method, status).Observe(float64(rw.size))

		// Record the httpResponseCount
		httpResponseCount.WithLabelValues(r.URL.Path, r.Method, status).Inc()

		// Record the httpResponseBucket
		httpResponseBucket.WithLabelValues(r.URL.Path, r.Method, status).Observe(duration)

		// Record the httpResponseIntervalCount with an example interval label
		interval := "default_interval" // Replace with actual interval logic if needed
		httpResponseIntervalCount.WithLabelValues(r.URL.Path, r.Method, status, interval).Inc()

		// Increment the error counter for error responses
		if rw.status >= 400 {
			errorCounter.WithLabelValues(r.URL.Path, r.Method, status).Inc()
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *responseWriter) WriteHeader(status int) {
	// Only write the status if it hasn't been written already
	if w.status == http.StatusOK {
		w.status = status
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *responseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}
