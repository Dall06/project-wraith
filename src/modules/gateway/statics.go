package gateway

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"project-wraith/src/pkg/logger"
)

type StaticsController interface {
	HelloHuman(ctx *fiber.Ctx) error
	LogReport(ctx *fiber.Ctx) error
	ResetPassword(ctx *fiber.Ctx) error
}

type staticsController struct {
	log      logger.Logger
	logsPath string
	basePath string
	version  string
}

func NewStaticsController(log logger.Logger, version, logsPath, basePath string) StaticsController {
	return &staticsController{
		log:      log,
		version:  version,
		logsPath: logsPath,
		basePath: basePath,
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
	sc.log.Info("Get Application Welcome")

	return ctx.Render("index", fiber.Map{
		"Version": sc.version,
	})
}

// LogReport
// @Summary Get Log Report
// @Description Returns the log report in the response rendered page.
// @Tags Static
// @Accept json
// @Produce html
// @Router /log [get]
// @Success 200 {string} string "HTML page with application version"
// @Failure 500 {object} error "Internal server error"
// @Security ApiKeyAuth
func (sc *staticsController) LogReport(ctx *fiber.Ctx) error {
	sc.log.Info("Get Log Report")

	infoLogsPath := fmt.Sprintf("%s/info.log", sc.logsPath)
	errLogsPath := fmt.Sprintf("%s/error.log", sc.logsPath)
	warnLogsPath := fmt.Sprintf("%s/warn.log", sc.logsPath)

	infoLogsContent, err := logger.ReadFile(infoLogsPath)
	if err != nil {
		return err
	}

	errLogsContent, err := logger.ReadFile(errLogsPath)
	if err != nil {
		return err
	}

	warnLogsContent, err := logger.ReadFile(warnLogsPath)
	if err != nil {
		return err
	}

	return ctx.Render("logs", fiber.Map{
		"Version":   sc.version,
		"InfoLogs":  infoLogsContent,
		"ErrorLogs": errLogsContent,
		"WarnLogs":  warnLogsContent,
	})
}

// ResetPassword
// @Summary Reset Password
// @Description Resets the user's password with the provided new password.
// @Tags Reset
// @Accept json
// @Produce json
// @Router /reset/modify [post]
// @Param request body Reset true "Reset request object"
// @Success 200 {object} map[string]string "Password reset successful"
// @Failure 400 {object} error "Failed to reset password"
// @Security ApiKeyAuth
func (sc *staticsController) ResetPassword(ctx *fiber.Ctx) error {
	token := utils.CopyString(ctx.Params("token"))

	return ctx.Render("reset", fiber.Map{
		"Token":     token,
		"Version":   sc.version,
		"ResetPath": fmt.Sprintf("%s/reset/modify", sc.basePath),
	})
}
