package core

import (
	"encoding/base64"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"project-wraith/src/pkg/guard"
	"project-wraith/src/pkg/link"
	"project-wraith/src/pkg/logger"
	"project-wraith/src/pkg/token"
	"time"
)

func Helmet() fiber.Handler {
	cfg := helmet.Config{
		CSPReportOnly: true,
	}

	return helmet.New(cfg)
}

func Compress() fiber.Handler {
	cfg := compress.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Method() != fiber.MethodGet
		},
		Level: compress.LevelBestSpeed, // 1
	}
	return compress.New(cfg)
}

func EncryptCookie(secret string) fiber.Handler {
	cfg := encryptcookie.Config{
		Key: base64.StdEncoding.EncodeToString([]byte(secret)),
	}
	return encryptcookie.New(cfg)
}

func ETag() fiber.Handler {
	cfg := etag.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Method() != fiber.MethodGet
		},
		Weak: true,
	}
	return etag.New(cfg)
}

func Recover() fiber.Handler {
	return recover.New()
}

func JwtWare(jwtSecret string, lookUp string) fiber.Handler {
	cfg := jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(jwtSecret),
		},
		TokenLookup: lookUp,
	}

	return jwtware.New(cfg)
}

func KeyAuth(apiKey string) fiber.Handler {
	cfg := keyauth.Config{
		KeyLookup: "header:x-access-token",
		Validator: func(c *fiber.Ctx, s string) (bool, error) {
			if apiKey != s {
				return false, nil
			}

			return true, nil
		},
	}
	return keyauth.New(cfg)
}

func CRSF() fiber.Handler {
	cfg := csrf.Config{
		Expiration: 15 * time.Minute,
	}

	return csrf.New(cfg)
}

func CORS() fiber.Handler {
	cfg := &cors.Config{
		AllowOrigins:  "*",
		AllowHeaders:  "Origin,Content-Type,Accept,X-Session-Token,X-Application-Key",
		AllowMethods:  "GET,POST,PUT,DELETE",
		ExposeHeaders: "Content-Length,Authorization",
		MaxAge:        5600,
	}
	return cors.New(*cfg)
}

func ResetAuth(jwtSecret string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkn := ctx.Get("X-Reset-Token")
		if tkn == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{Message: "invalid token"})
		}

		valid, err := token.ValidateJwtToken(tkn, jwtSecret, nil)
		if err != nil {
			return err
		}

		if !valid {
			return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{Message: "invalid token"})
		}

		return ctx.Next()
	}
}

func ManticoreSight(manticore guard.Manticore, log logger.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		cred := guard.Credentials{
			Username: ctx.FormValue("username"),
			Password: ctx.FormValue("password"),
		}

		err := manticore.StingAndProwl(cred)
		if err != nil {
			log.Error("failed to validate token: %v", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{Message: err.Error()})
		}

		return ctx.Next()
	}
}
