package gateway

import (
	"github.com/gofiber/fiber/v2"
	"project-wraith/src/modules/rules"
	"project-wraith/src/pkg/alchemy"
	"project-wraith/src/pkg/link"
	"project-wraith/src/pkg/logger"
	"project-wraith/src/pkg/token"
	"time"
)

type UserController interface {
	Register(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	Edit(ctx *fiber.Ctx) error
	Remove(ctx *fiber.Ctx) error
}

type userController struct {
	log                logger.Logger
	rules              rules.UserRule
	jwtSecret          string
	encryptResponse    bool
	responseSecret     string
	cookiesMinutesLife time.Duration
}

func NewUserController(
	log logger.Logger,
	rules rules.UserRule,
	jwtSecret string,
	encryptResponse bool,
	responseSecret string,
	cookiesMinutesLife int,
) UserController {
	return &userController{
		log:                log,
		rules:              rules,
		jwtSecret:          jwtSecret,
		encryptResponse:    encryptResponse,
		responseSecret:     responseSecret,
		cookiesMinutesLife: time.Duration(cookiesMinutesLife) * time.Minute,
	}
}

// Login
// @Summary User login
// @Description Authenticates a user and generates a session token if the credentials are valid.
// @Tags User
// @Accept json
// @Produce json
// @Router /user/login [post]
// @Param request body User true "User login credentials"
// @Success 200 {object} map[string]string "Login successful with session token"
// @Failure 400 {object} error "Failed to parse request or invalid credentials"
// @Failure 401 {object} error "Unauthorized access"
// @Failure 500 {object} error "Internal server error"
// @Security ApiKeyAuth
func (uc userController) Login(ctx *fiber.Ctx) error {
	req := User{}
	if err := ctx.BodyParser(&req); err != nil {
		uc.log.Error("failed to parse request: %v", err)
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

	res, err := uc.rules.Login(actor)
	if err != nil {
		uc.log.Error("failed to login: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(link.Response{
			Message: err.Error(),
		})
	}

	userSession, err := token.CreateJwtToken(
		uc.jwtSecret, uc.cookiesMinutesLife, res)
	if err != nil {
		uc.log.Error("failed to create token token: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(link.Response{
			Message: err.Error(),
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "user_session",
		Value:    userSession,
		Expires:  time.Now().Add(uc.cookiesMinutesLife),
		HTTPOnly: true,
		Secure:   true,
	})

	uc.log.Info("action done: login successful")
	return ctx.Status(fiber.StatusOK).JSON(link.Response{
		Message: "login successful",
	})
}

// Exit
// @Summary User logout
// @Description Logs out the user by expiring the session token.
// @Tags User
// @Accept json
// @Produce json
// @Router /user/logout [post]
// @Success 200 {object} map[string]string "Logout successful"
// @Failure 401 {object} error "No session found"
// @Failure 500 {object} error "Failed to expire session"
// @Security ApiKeyAuth
func (uc userController) Exit(ctx *fiber.Ctx) error {
	userSession := ctx.Cookies("user_session")

	if userSession == "" {
		uc.log.Error("no session found")
		return ctx.Status(fiber.StatusUnauthorized).JSON(link.Response{
			Message: "no session found",
		})
	}

	expiredToken, err := token.ExpireJwtToken(uc.jwtSecret, time.Hour, nil)
	if err != nil {
		uc.log.Error("failed to expire session: %v", err)
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

	uc.log.Info("action successful: logout user")
	return ctx.Status(fiber.StatusOK).JSON(link.Response{
		Message: "logout successful",
	})
}

// Register
// @Summary User registration
// @Description Registers a new user with the provided details.
// @Tags User
// @Accept json
// @Produce json
// @Router /user/register [post]
// @Param request body User true "User registration details"
// @Success 200 {object} map[string]string "Registration successful"
// @Failure 400 {object} error "Failed to parse request or registration error"
// @Security ApiKeyAuth
func (uc userController) Register(ctx *fiber.Ctx) error {
	req := User{}
	if err := ctx.BodyParser(&req); err != nil {
		uc.log.Error("failed to parse request: %v", err)
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

	res, err := uc.rules.Register(actor)
	if err != nil {
		uc.log.Error("failed to register: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(link.Response{
			Message: err.Error(),
		})
	}

	uc.log.Info("action done: register successful %v", res)
	return ctx.Status(fiber.StatusOK).JSON(link.Response{
		Message: "register successful",
	})
}

// Get
// @Summary Get user details
// @Description Retrieves user details based on the provided user ID.
// @Tags User
// @Accept json
// @Produce json
// @Router /user/{id} [get]
// @Param id path string true "User ID"
// @Success 200 {object} User "User details"
// @Failure 400 {object} error "Invalid ID or request"
// @Failure 404 {object} error "User not found"
// @Security ApiKeyAuth
func (uc userController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		uc.log.Error("parameter not found: {key: id, value: %v}", id)
		return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{
			Message: "id is required",
		})
	}

	actor := rules.User{ID: id}

	user, err := uc.rules.Get(actor)
	if err != nil {
		uc.log.Warn("failed to get user: %v", err)
		return ctx.Status(fiber.StatusOK).JSON(link.Response{
			Message: err.Error(),
		})
	}

	if user == nil {
		uc.log.Warn("empty data: user not found")
		return ctx.Status(fiber.StatusOK).JSON(link.Response{
			Message: "user not found",
		})
	}

	res := User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
		Phone:    user.Phone,
	}

	uc.log.Info("action done: get user")

	if uc.encryptResponse {
		structString, err := alchemy.StructIntoString(&res, uc.responseSecret)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(structString)
	}

	return ctx.Status(fiber.StatusOK).JSON(link.Response{
		Content: res,
	})
}

// Edit
// @Summary Edit user details
// @Description Updates the user details with the provided information.
// @Tags User
// @Accept json
// @Produce json
// @Router /user/edit [put]
// @Param request body User true "Updated user details"
// @Success 200 {object} map[string]string "User details updated successfully"
// @Failure 400 {object} error "Failed to parse request or update error"
// @Security ApiKeyAuth
func (uc userController) Edit(ctx *fiber.Ctx) error {
	req := User{}
	if err := ctx.BodyParser(&req); err != nil {
		uc.log.Error("failed to parse request: %v", err)
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

	err := uc.rules.Edit(actor)
	if err != nil {
		uc.log.Error("failed to edit user: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(link.Response{
			Message: err.Error(),
		})
	}

	uc.log.Info("action successful: edit user")
	return ctx.Status(fiber.StatusOK).JSON(link.Response{
		Message: "edit successful",
	})
}

// Remove
// @Summary Remove user
// @Description Deletes a user based on the provided details.
// @Tags User
// @Accept json
// @Produce json
// @Router /user/remove [delete]
// @Param request body User true "User details for removal"
// @Success 200 {object} map[string]string "User removed successfully"
// @Failure 400 {object} error "Failed to parse request or removal error"
// @Security ApiKeyAuth
func (uc userController) Remove(ctx *fiber.Ctx) error {
	req := User{}
	if err := ctx.BodyParser(&req); err != nil {
		uc.log.Error("failed to parse request: %v", err)
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

	err := uc.rules.Remove(actor)
	if err != nil {
		uc.log.Error("failed to remove user: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(link.Response{
			Message: err.Error(),
		})
	}

	uc.log.Info("action successful: remove user")
	return ctx.Status(fiber.StatusOK).JSON(link.Response{
		Message: "remove successful",
	})
}
