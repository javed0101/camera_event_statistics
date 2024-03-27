package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func RegisterPrometheusMetrics() {
	prometheus.MustRegister(totalRequests, latency, camEventCounter)
}

func StatusCodeMetrics(c *fiber.Ctx) {

	start := time.Now()
	// next := c.Next()
	elapsed := time.Since(start).Seconds()

	latency.WithLabelValues(
		c.Route().Method,
		c.Route().Path,
	).Observe(elapsed)

	statusCode := c.Response().StatusCode()

	if statusCode >= 200 && statusCode < 300 {
		totalRequests.WithLabelValues("2xx", "/api/v1/camera/event/stats", "GET").Inc()
	} else if statusCode >= 400 && statusCode < 500 {
		totalRequests.WithLabelValues("4xx", "/api/v1/camera/event/stats", "GET").Inc()
	} else if statusCode >= 500 && statusCode < 599 {
		totalRequests.WithLabelValues("5xx", "/api/v1/camera/event/stats", "GET").Inc()
	}
}

func CameraEventMetrics(c *fiber.Ctx) error {

	cameraID := c.Query("cameraID")
	eventType := c.Query("eventType")

	camEventCounter.WithLabelValues(cameraID + eventType).Inc()

	return c.Next()
}

func LatencyMetrics(c *fiber.Ctx) error {
	start := time.Now()
	next := c.Next()
	elapsed := time.Since(start).Seconds()

	latency.WithLabelValues(
		c.Route().Method,
		c.Route().Path,
	).Observe(elapsed)

	return next
}

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"status_code", "path", "method"},
	)
)

var (
	camEventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "camera_event_counter",
			Help: "Metic to track camera events",
		},
		[]string{"cam_event"},
	)
)

// var latency = prometheus.NewHistogramVec(
// 	prometheus.HistogramOpts{
// 		Namespace: "api",
// 		Name:      "latency_seconds",
// 		Help:      "Latency distributions.",
// 		Buckets:   []float64{0.1, 0.2, 0.5},
// 	},
// 	[]string{"method", "path"},
// )

var latency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "api",
		Name:       "latency_seconds",
		Help:       "Latency distributions.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"method", "path"},
)
