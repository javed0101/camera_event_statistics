package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

// func RegisterPrometheusMetrics() {
// 	prometheus.MustRegister(latency, totalRequests, tps)
// }

// func RecordRequestLatency(c *fiber.Ctx) error {
// 	start := time.Now()
// 	next := c.Next()
// 	elapsed := time.Since(start).Seconds()

// 	latency.WithLabelValues(
// 		c.Route().Method,
// 		c.Route().Path,
// 	).Observe(elapsed)

// 	return next
// }

func CameraEventMiddleware(c *fiber.Ctx) error {

	cameraID := c.Query("cameraID")
	eventType := c.Query("eventType")

	camEventCounter.WithLabelValues(cameraID + eventType).Inc()

	return c.Next()
}

var latency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "api",
		Name:      "latency_seconds",
		Help:      "Latency distributions.",
		Buckets:   []float64{0.99, 0.90, 0.50},
	},
	[]string{"method", "path"},
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
