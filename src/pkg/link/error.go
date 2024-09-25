package link

import "github.com/gofiber/fiber/v2"

func Error(ctx *fiber.Ctx, err error) error {

	return ctx.Status(fiber.StatusBadRequest).JSON(Response{Message: err.Error()})
}
