package render

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
	"time"
	"void-studio.net/fiesta/fapi/auth"
	"void-studio.net/fiesta/fapi/config"
)

func Login(ctx *fiber.Ctx) error {
	if ctx.Query("perms") == "none" {
		return ctx.Render("login", nil)
	}

	if code := ctx.Query("code"); code != "" {
		if discordToken := auth.GetUserToken(code); discordToken != "" && CreateCookie(ctx, discordToken) {
			return ctx.RedirectToRoute("main", fiber.Map{}, fasthttp.StatusOK)
		}

		return ctx.RedirectToRoute("login", fiber.Map{"perms": "none"}, fasthttp.StatusUnauthorized)
	}

	return ctx.Redirect(config.Config.Discord.AuthURI)
}

func CreateCookie(ctx *fiber.Ctx, discordToken string) bool {
	ctx.ClearCookie("Authorization")
	accessed := auth.GetAccessRole(discordToken)

	if accessed == nil {
		return false
	}

	claims := jwt.MapClaims{
		config.AccessServersKey: accessed,
		config.DiscordTokenKey:  discordToken,
		config.NeedCheckKey:     time.Now().Add(time.Hour * 24).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signedToken, err := jwtToken.SignedString(config.Config.General.Secret)
	if err != nil {
		log.Errorf("failed to sign token: %v", err)
		return false
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    signedToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   !config.Config.General.DevMode,
		SameSite: "Strict",
	})

	return true
}
