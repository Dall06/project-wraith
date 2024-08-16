package core

import (
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
	"project-wraith/src/modules/rules"
	"project-wraith/src/pkg/logger"
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
		Key: secret,
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

func ResetPassword(resetRule rules.ResetRule, log logger.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tkn := ctx.Get("X-Reset-Token")
		if tkn == "" {
			log.Error("parameter not found: {key: X-Reset-Token, value: %v}", tkn)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid token"})
		}

		reset := rules.Reset{
			Token: tkn,
		}

		err := resetRule.Validate(reset)
		if err != nil {
			log.Error("failed to validate token: %v", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return nil
	}
}
