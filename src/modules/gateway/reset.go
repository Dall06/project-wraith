package gateway

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"project-wraith/src/modules/rules"
	"project-wraith/src/pkg/logger"
	"project-wraith/src/pkg/mail"
	"project-wraith/src/pkg/sms"
	"time"
)

type ResetController interface {
	Start(ctx *fiber.Ctx) error
	Modify(ctx *fiber.Ctx) error
}

type resetController struct {
	log              logger.Logger
	reset            rules.ResetRule
	user             rules.UserRule
	jwtSecret        string
	cookieExpiration time.Duration
	mailer           mail.Mail
	smsSender        sms.Twilio
	webUrl           string
}

func NewResetController(
	log logger.Logger,
	reset rules.ResetRule,
	user rules.UserRule,
	jwtSecret string,
	cookieExpiration int,
	mailer mail.Mail,
	smsSender sms.Twilio,
	webUrl string,
) ResetController {
	return &resetController{
		log:              log,
		reset:            reset,
		user:             user,
		jwtSecret:        jwtSecret,
		cookieExpiration: time.Duration(cookieExpiration) * time.Hour,
		mailer:           mailer,
		smsSender:        smsSender,
		webUrl:           webUrl,
	}
}

// Start
// @Summary Start password reset
// @Description Initiates the password reset process by sending a reset link to the user's email and SMS.
// @Tags Reset
// @Accept json
// @Produce json
// @Router /reset/start [post]
// @Param request body Reset true "Reset request object"
// @Success 202 {object} map[string]string "Successfully sent password reset link"
// @Failure 400 {object} error "Failed to parse request or send notifications"
// @Security ApiKeyAuth
func (rc resetController) Start(ctx *fiber.Ctx) error {
	req := Reset{}
	if err := ctx.BodyParser(&req); err != nil {
		rc.log.Error("failed to parse request: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse request",
		})
	}

	model := rules.Reset{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	entity, err := rc.reset.Start(model)
	if err != nil {
		rc.log.Error("failed to get user: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resetWebUrl := fmt.Sprintf(rc.webUrl, entity.Token)

	bindStruct := struct {
		Username string
		ResetUrl string
	}{
		Username: entity.Username,
		ResetUrl: resetWebUrl,
	}

	err = rc.mailer.Send(
		"./templates/email.html",
		bindStruct,
		"Reset Password",
		[]string{entity.Email})
	if err != nil {
		rc.log.Error("failed to send mail: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	res, err := rc.smsSender.SendSMSTwilio(
		entity.Phone, true, resetWebUrl)
	if err != nil {
		rc.log.Error("failed to send sms: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if res == "" {
		rc.log.Error("failed to send sms: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to send sms",
		})
	}

	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "password reset link sent",
		"token":   entity.Token,
	})

}

// Modify
// @Summary Modify password
// @Description Resets the user's password with the provided new password.
// @Tags Reset
// @Accept json
// @Produce json
// @Router /reset/modify [post]
// @Param request body Reset true "Reset request object"
// @Success 200 {object} map[string]string "Password reset successful"
// @Failure 400 {object} error "Failed to reset password"
// @Security ApiKeyAuth
func (rc resetController) Modify(ctx *fiber.Ctx) error {
	var req Reset
	if err := ctx.BodyParser(&req); err != nil {
		rc.log.Error("failed to parse request: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse request",
		})
	}

	model := rules.User{
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.NewPassword,
	}

	err := rc.user.Edit(model)
	if err != nil {
		rc.log.Error("failed to reset password: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "password reset successful",
	})
}
