package api

import (
	"github.com/gofiber/fiber/v2"
	"void-studio.net/fiesta/fapi/config"
)

func Search(ctx *fiber.Ctx) error {
	if ctx.Locals(config.AuthKey) != nil {

	}

	return nil
}
