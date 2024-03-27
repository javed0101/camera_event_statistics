package routes

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/javed0101/cameraevents/internal/core/handlers"
	"github.com/javed0101/cameraevents/pkg/middlewares"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRoutes(app *fiber.App) {

	middlewares.FiberMiddlewares(app)
	handlers.RegisterPrometheusMetrics()

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	app.Get("/", handlers.RootHandler)
	app.Post("/", handlers.RootHandler)

	route := app.Group("/api/v1")

	route.Get("/health", handlers.HealthHandler)

	route.Get("/camera/event/stats", handlers.CameraEventMetrics, handlers.CameraEventHandler)

	// go handlers.ResetTotalRequestsPeriodically()

	app.Use(handlers.NoRouteHandler)
}
