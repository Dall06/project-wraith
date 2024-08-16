package gateway

import (
	"github.com/gofiber/fiber/v2"
	"project-wraith/src/pkg/logger"
)

type StaticsController interface {
	HelloHuman(ctx *fiber.Ctx) error
}

type staticsController struct {
	log     logger.Logger
	version string
}

func NewStaticsController(log logger.Logger, version string) StaticsController {
	return &staticsController{
		log:     log,
		version: version,
	}
}

// HelloHuman
// @Summary Get Application Version
// @Description Returns the application version in the response rendered page.
// @Tags Static
// @Accept json
// @Produce html
// @Router /hello [get]
// @Success 200 {string} string "HTML page with application version"
// @Failure 500 {object} error "Internal server error"
// @Security ApiKeyAuth
func (sc *staticsController) HelloHuman(ctx *fiber.Ctx) error {
	return ctx.Render("index", fiber.Map{
		"Version": sc.version,
	})
}
