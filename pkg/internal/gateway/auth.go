package gateway

import (
	"github.com/gofiber/fiber/v2"
	"project-wraith/pkg/internal/rules"
	"project-wraith/pkg/modules/link"
	"project-wraith/pkg/modules/logger"
	"project-wraith/pkg/modules/token"
	"time"
)

type AuthController interface {
	Login(ctx *fiber.Ctx) error
	Exit(ctx *fiber.Ctx) error
}

type authController struct {
	log                logger.Logger
	rules              rules.UserRule
	jwtSecret          string
	cookiesMinutesLife time.Duration
}

func NewAuthController(
	log logger.Logger,
	rules rules.UserRule,
	jwtSecret string,
	cookiesMinutesLife int,
) AuthController {
	return &authController{
		log:                log,
		rules:              rules,
		jwtSecret:          jwtSecret,
		cookiesMinutesLife: time.Duration(cookiesMinutesLife) * time.Minute,
	}
}

// Login
// @Summary Auth login
// @Description Authenticates a user and generates a session token if the credentials are valid.
// @Tags Auth
// @Accept json
// @Produce json
// @Router /user/login [post]
// @Param request body Auth true "Auth login credentials"
// @Success 200 {object} map[string]string "Login successful with session token"
// @Failure 400 {object} error "Failed to parse request or invalid credentials"
// @Failure 401 {object} error "Unauthorized access"
// @Failure 500 {object} error "Internal server error"
// @Security ApiKeyAuth
func (ac authController) Login(ctx *fiber.Ctx) error {
	req := User{}
	if err := ctx.BodyParser(&req); err != nil {
		ac.log.Error("failed to parse request: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{
			Message: "failed to parse request",
		})
	}

	actor := rules.User{
		ID:       req.ID,
		Username: req.Username,
		Email:    req.Email,
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
	}

	res, err := ac.rules.Login(actor)
	if err != nil {
		ac.log.Error("failed to login: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(link.Response{
			Message: err.Error(),
		})
	}

	userSession, err := token.CreateJwtToken(
		ac.jwtSecret, ac.cookiesMinutesLife, res)
	if err != nil {
		ac.log.Error("failed to create token token: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(link.Response{
			Message: err.Error(),
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "user_session",
		Value:    userSession,
		Expires:  time.Now().Add(ac.cookiesMinutesLife),
		HTTPOnly: true,
		Secure:   true,
	})

	ac.log.Info("action done: login successful")
	return ctx.Status(fiber.StatusOK).JSON(link.Response{
		Message: "login successful",
	})
}

// Exit
// @Summary Auth logout
// @Description Logs out the user by expiring the session token.
// @Tags Auth
// @Accept json
// @Produce json
// @Router /user/logout [post]
// @Success 200 {object} map[string]string "Logout successful"
// @Failure 401 {object} error "No session found"
// @Failure 500 {object} error "Failed to expire session"
// @Security ApiKeyAuth
func (ac authController) Exit(ctx *fiber.Ctx) error {
	userSession := ctx.Cookies("user_session")

	if userSession == "" {
		ac.log.Error("no session found")
		return ctx.Status(fiber.StatusUnauthorized).JSON(link.Response{
			Message: "no session found",
		})
	}

	expiredToken, err := token.ExpireJwtToken(ac.jwtSecret, time.Hour, nil)
	if err != nil {
		ac.log.Error("failed to expire session: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(link.Response{
			Message: "failed to expire session",
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "user_session",
		Value:    expiredToken,
		Expires:  time.Unix(0, 0), // Set to a time in the past
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	ac.log.Info("action successful: logout user")
	return ctx.Status(fiber.StatusOK).JSON(link.Response{
		Message: "logout successful",
	})
}
