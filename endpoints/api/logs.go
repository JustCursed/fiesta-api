package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Logs(ctx *fiber.Ctx) error {
	//if ctx.Locals(config.AuthKey).(bool) {
	//
	//}

	log.Info(ctx.Queries())

	return ctx.JSON(fiber.Map{"logs": []string{"<UNK>", "<UNK>", "<UNK>", "<UNK>", "<UNK>", "<UNK>", "<UNK>", "<UNK>"}})
}
