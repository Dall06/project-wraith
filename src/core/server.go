package core

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/swagger"
	"project-wraith/src/modules/gateway"
	"project-wraith/src/pkg/guard"
	"project-wraith/src/pkg/logger"
)

func Middleware(
	app *fiber.App,
	log logger.Logger,
	paths map[string]string,
	serverApiKey,
	jwtSecret,
	cookiesSecret string,
	manticore guard.Manticore) {

	app.Use(CORS())
	app.Use(Compress())
	app.Use(ETag())
	app.Use(Helmet())
	app.Use(Recover())
	app.Use(CRSF())
	app.Use(EncryptCookie(cookiesSecret))

	for key, path := range paths {
		if key != "hello" {
			app.Use(path, KeyAuth(serverApiKey))
		}

		switch key {
		case "user":
			app.Use(path, JwtWare(jwtSecret, "cookie:user_session"))
		case "reset":
			app.Use(fmt.Sprintf("%s/form", path), ResetAuth(jwtSecret))
		case "logs":
			app.Use(path, ManticoreSight(manticore, log))
		case "metrics":
			app.Use(path, ManticoreSight(manticore, log))
		default:
			continue
		}
	}
}

func EnRoute(
	app *fiber.App,
	paths map[string]string,
	user gateway.UserController,
	auth gateway.AuthController,
	reset gateway.ResetController,
	statics gateway.StaticsController) {

	for key, path := range paths {
		switch key {
		case "hello":
			app.Get(path, statics.HelloHuman)
		case "logs":
			app.Get(path, statics.LogReport)
		case "metrics":
			app.Get(paths["metrics"], monitor.New(monitor.Config{
				Title: "project-wraith metrics",
			}))
		case "swagger":
			app.Get(path, swagger.HandlerDefault)
		case "auth":
			authGroup := app.Group(path)
			authGroup.Post("/login", auth.Login)
			authGroup.Put("/exit", auth.Exit)
		case "reset":
			passResetGroup := app.Group(path)
			passResetGroup.Post("/init", reset.Start)
			passResetGroup.Post("/modify", reset.Modify)
		case "user":
			usersGroup := app.Group(path)
			usersGroup.Post("/register", user.Register)
			usersGroup.Get("/detail/:id", user.Get)
			usersGroup.Put("/edit", user.Edit)
			usersGroup.Delete("/remove", user.Remove)
		default:
			continue
		}
	}
}
