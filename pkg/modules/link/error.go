package link

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

func Error(ctx *fiber.Ctx, err error) error {
	var errorFiber *fiber.Error
	ok := errors.As(err, &errorFiber)
	if ok {
		return ctx.Status(errorFiber.Code).JSON(Response{Message: errorFiber.Error()})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(Response{Message: err.Error()})
}
