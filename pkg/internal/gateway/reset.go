package gateway

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"project-wraith/pkg/internal/rules"
	"project-wraith/pkg/modules/link"
	"project-wraith/pkg/modules/logger"
	"project-wraith/pkg/modules/mail"
	"project-wraith/pkg/modules/sms"
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
		return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{
			Message: "failed to parse request",
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
		return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{
			Message: err.Error(),
		})
	}

	bindStruct := struct {
		Username   string
		Email      string
		Phone      string
		ResetToken string
	}{
		Username:   entity.Username,
		Email:      entity.Email,
		Phone:      entity.Phone,
		ResetToken: entity.Token,
	}

	err = rc.mailer.Send(
		"./templates/email.html",
		bindStruct,
		"Reset Password",
		[]string{entity.Email})
	if err != nil {
		rc.log.Error("failed to send mail: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{
			Message: err.Error(),
		})
	}

	resetWebUrl := fmt.Sprintf("%s/%s", rc.webUrl, entity.Token)

	res, err := rc.smsSender.SendSMSTwilio(
		entity.Phone, true, resetWebUrl)
	if err != nil {
		rc.log.Error("failed to send sms: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{
			Message: err.Error(),
		})
	}
	if res == "" {
		rc.log.Error("failed to send sms: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{
			Message: "failed to send sms",
		})
	}

	return ctx.Status(fiber.StatusAccepted).JSON(link.Response{
		Message: "password reset link sent",
		Content: Outcome{
			Token:    entity.Token,
			ResetUrl: resetWebUrl,
		},
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
	tkn := ctx.Get("X-Reset-Token")
	if tkn == "" {
		log.Error("parameter not found: {key: X-Reset-Token, value: %v}", tkn)
		return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{Message: "invalid token"})
	}

	reset := rules.Reset{
		Token: tkn,
	}
	res, err := rc.reset.Validate(reset)
	if err != nil {
		log.Error("failed to validate: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{Message: err.Error()})
	}

	req := Reset{}
	if err := ctx.BodyParser(&req); err != nil {
		log.Error("failed to parse request: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{Message: "failed to parse request"})
	}

	model := rules.User{
		ID:       res.ID,
		Password: req.NewPassword,
	}
	err = rc.user.Edit(model)
	if err != nil {
		rc.log.Error("failed to reset password: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(link.Response{
		Message: "password reset successful",
	})
}
