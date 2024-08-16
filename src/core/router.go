package core

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"project-wraith/src/modules/gateway"
)

func Middleware(app *fiber.App, basePath, serverApiKey, jwtSecret, cookiesSecret string) {
	// Define paths
	userPath := fmt.Sprintf("%s/user", basePath)
	authPath := fmt.Sprintf("%s/auth", basePath)
	resetPath := fmt.Sprintf("%s/reset", basePath)

	// Global Middleware
	app.Use(CORS())
	app.Use(Helmet())
	app.Use(Compress())
	app.Use(ETag())
	app.Use(CRSF())
	app.Use(Recover())

	// Middleware for Authentication Paths
	app.Use(authPath, KeyAuth(serverApiKey))
	app.Use(authPath, EncryptCookie(cookiesSecret))

	// Middleware for User Paths
	app.Use(userPath, KeyAuth(serverApiKey))
	app.Use(userPath, EncryptCookie(cookiesSecret))
	app.Use(userPath, JwtWare(jwtSecret, "cookie:user_session"))

	// Middleware for Reset Paths
	app.Use(resetPath, KeyAuth(serverApiKey))
	app.Use(resetPath, EncryptCookie(cookiesSecret))
}

func EnRoute(
	app *fiber.App,
	basePath string,
	user gateway.UserController,
	reset gateway.ResetController,
	statics gateway.StaticsController) {
	app.Get(fmt.Sprintf("%s/hello", basePath), statics.HelloHuman)

	app.Get(fmt.Sprintf("%s/swagger/*", basePath), swagger.HandlerDefault)

	authGroup := app.Group(fmt.Sprintf("%s/auth", basePath))
	authGroup.Post("/login", user.Login)
	authGroup.Put("/exit", user.Exit)

	passResetGroup := app.Group(fmt.Sprintf("%s/reset", basePath))
	passResetGroup.Post("/init", reset.Start)
	passResetGroup.Post("/modify", reset.Modify)

	usersGroup := app.Group(fmt.Sprintf("%s/user", basePath))
	usersGroup.Post("/register", user.Register)
	usersGroup.Get("/detail/:id", user.Get)
	usersGroup.Put("/edit", user.Edit)
	usersGroup.Delete("/remove", user.Remove)
}
