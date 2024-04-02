package http

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/javed0101/cameraevents/cmd/routes"
	configmanager "github.com/javed0101/cameraevents/config"
	redis "github.com/javed0101/cameraevents/internal/sources/redis"
	"github.com/javed0101/cameraevents/processor"
)

func InitAPIServer() {

	appEnv := os.Getenv("APP_ENV")
	log.Printf("App is running in %s environment", appEnv)

	configmanager.InitConfig(appEnv)
	config := configmanager.GetConfig()

	redis.InitRedis()

	// Init pulsar client and workers
	// Sleeping for 5s because pulsar is taking some time to come up, depends_on is not working in docker-compose
	time.Sleep(time.Second * 5)
	go processor.JobProcessor(*config.Pulsar.Topic.CameraEvent).ProcessJob()

	app := fiber.New()

	routes.SetupRoutes(app)

	err := app.Listen(config.App.Port)
	if err != nil {
		log.Fatalf("Error initializing server. Exiting with error: %s", err.Error())
	}
}
