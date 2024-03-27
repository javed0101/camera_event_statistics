package handlers

import (
	"github.com/gofiber/fiber/v2"
	eenerror "github.com/javed0101/cameraevents/pkg/utils"
)

func RootHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

func HealthHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "healthy",
	})
}

func NoRouteHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": eenerror.ErrorInvalidRequest,
	})
}
