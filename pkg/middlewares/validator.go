package middlewares

import "github.com/gofiber/fiber/v2"

func CommonValidator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		contentType := c.Get("Content-Type")
		if contentType != "application/json" {
			return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
				"error": "Unsupported Content Type",
			})
		}
		return c.Next()
	}
}
