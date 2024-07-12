package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler", "method", "status"},
	)

	requestsByStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by status code.",
		},
		[]string{"handler", "method", "status"},
	)

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
)

func init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestsByStatus)
	prometheus.MustRegister(postCounter)
	prometheus.MustRegister(getCounter)
	prometheus.MustRegister(patchCounter)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture the status code
		rw := &responseWriter{w, http.StatusOK}

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
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func MonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pass the request through the handler chain
		next.ServeHTTP(w, r)

		// Get the status code of the response
		status := strconv.Itoa(w.(*responseWriter).status)

		switch r.Method {
		case http.MethodPost:
			postCounter.WithLabelValues(status).Inc()
		case http.MethodGet:
			getCounter.WithLabelValues(status).Inc()
		case http.MethodPatch:
			patchCounter.WithLabelValues(status).Inc()
		}
	})
}
