package render

import (
	"github.com/gofiber/fiber/v2"
)

func MainPage(ctx *fiber.Ctx) error {
	return ctx.Render("index", fiber.Map{"dates": []string{"2021-03-10", "2022-03-10", "2023-03-10", "2024-03-10", "2025-03-10"}, "servers": []string{"HiTech", "MagicCraft", "GregTech", "ParasiteTech"}})
}
