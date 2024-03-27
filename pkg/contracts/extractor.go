package contracts

import "github.com/gofiber/fiber/v2"

type Extractor interface {
	ExtractFromRequest(*fiber.Ctx) *error
}
