package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequests = prometheus.NewCounterVec( 
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code, full path, and HTTP method.",
		}, 
		[]string{"method", "path", "status"}, 
	)

	httpLatency = prometheus.NewHistogramVec( 
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "The amount of time for each request to be completed.",
			Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1, 2, 5},
		}, 
		[]string{"method", "path"}, 
	)

	authResults = prometheus.NewCounterVec( 
		prometheus.CounterOpts{
			Name: "auth_results_total",
			Help: "How many HTTP authentication requests processed, partitioned by either success or failure.",
		}, 
		[]string{"result"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests, httpLatency, authResults)
}

func Metrics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		httpRequests.WithLabelValues(
			ctx.Request.Method,
			ctx.FullPath(),
			strconv.Itoa(ctx.Writer.Status()),
		).Inc()

		httpLatency.WithLabelValues(
			ctx.Request.Method,
			ctx.FullPath(),
		).Observe(time.Since(start).Seconds())

		if authResult := ctx.GetString("auth_result"); authResult != "" {
			authResults.WithLabelValues(
				authResult,
			).Inc()
		}

	}
}
