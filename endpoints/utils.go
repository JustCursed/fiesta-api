package endpoints

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
//	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"void-studio.net/fiesta/fapi/config"
	"void-studio.net/fiesta/fapi/endpoints/api"
	"void-studio.net/fiesta/fapi/endpoints/render"
)

func RegisterEndpoints(app *fiber.App) {
	app.Get("/", render.MainPage).Name("main")
	app.Get("/login", render.Login).Name("login")
	app.Get("/api/logs/:server/:type", api.Logs)
	app.Get("/api/search/:server/:type", api.Search)
}

func RegisterMiddlewares(app *fiber.App) {
	app.Use(jwtware.New(
		jwtware.Config{
			ContextKey:     "jwtToken",
			TokenLookup:    "cookie:Authorization",
			ErrorHandler:   func(ctx *fiber.Ctx, _ error) error { return ctx.Next() },
			SuccessHandler: func(ctx *fiber.Ctx) error { return ctx.Next() },
			SigningKey: jwtware.SigningKey{
				JWTAlg: jwtware.ES256,
				Key:    config.Config.General.Secret.Public(),
			},
		},
	))

	app.Use(checkAuthWare)
}

func checkAuthWare(ctx *fiber.Ctx) error {
//	log.Info(ctx.OriginalURL())
//	log.Info(ctx.BaseURL())
	if token := ctx.Locals("jwtToken"); token != nil {
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)

		if time.Now().Unix() > claims[config.NeedCheckKey].(int64) {
			if !render.CreateCookie(ctx, claims[config.DiscordTokenKey].(string)) {
				ctx.Locals(config.AuthKey, false)
				return ctx.Next()
			}
		}

		ctx.Locals(config.AuthKey, true)
		return ctx.Next()
	}

	ctx.Locals(config.AuthKey, false)
	return ctx.Next()
}
