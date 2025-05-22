package main

import (
//	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/template/html/v2"
	"void-studio.net/fiesta/fapi/config"
	"void-studio.net/fiesta/fapi/endpoints"
)

func main() {
	app := fiber.New(fiber.Config{
		Views:          html.New("./views", ".html"),
		Prefork:        true,
		GETOnly:        true,
//		ReadBufferSize: 2048,
//		JSONEncoder:    json.Marshal,
//		JSONDecoder:    json.Unmarshal,
	})

	endpoints.RegisterMiddlewares(app)
	endpoints.RegisterEndpoints(app)

	app.Static("/static", "./views/public")

	err := app.Listen(config.Config.General.Address)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
}
