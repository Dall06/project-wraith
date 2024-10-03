package gateway_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"project-wraith/pkg/internal/gateway"
	"project-wraith/pkg/internal/rules"
	"project-wraith/pkg/modules/logger"
	"testing"
)

func TestAuth(test *testing.T) {
	logMock := &logger.MockLogger{}
	ruleMock := &rules.MockUserRule{}

	logMock.On("Initialize").Return(nil)
	logMock.On("Info", mock.Anything).Return(nil)  // Mock the Info method
	logMock.On("Error", mock.Anything).Return(nil) // Ensure Error is mocked correctly
	logMock.On("Warn", mock.Anything).Return(nil)  // Ensure Warn is mocked correctly

	authCtrl := gateway.NewAuthController(
		logMock,
		ruleMock,
		"super_secret_key",
		60,
	)

	tests := []struct {
		name   string
		action string
		method string
		input  gateway.User
	}{
		{
			name:   "Test Login",
			action: "login",
			method: "POST",
			input: gateway.User{
				ID:       "1",
				Username: "testuser",
				Email:    "test@example.com",
				Name:     "Test User",
				Phone:    "1234567890",
				Password: "securepassword",
			},
		},
		{
			name:   "Test Exit",
			action: "exit",
			method: "PUT",
			input:  gateway.User{}, // No input required for Exit
		},
	}

	for _, tc := range tests {
		test.Run(tc.name, func(t *testing.T) {
			app := fiber.New()

			switch tc.action {
			case "login":
				// Mock the Login method to return the expected response
				ruleMock.On("Login", mock.Anything).Return(&rules.User{ID: "1"}, nil).Once()
				app.Post("/login", authCtrl.Login)

				// Convert input to JSON for the request body
				inputBody, _ := json.Marshal(tc.input)
				req := httptest.NewRequest(tc.method, fmt.Sprintf("/%s", tc.action), bytes.NewBuffer(inputBody))
				req.Header.Set("Content-Type", "application/json") // Set Content-Type header for JSON

				resp, err := app.Test(req, -1)
				if err != nil {
					t.Fatalf("Fiber test error: %v", err)
				}

				if resp.StatusCode != fiber.StatusOK {
					t.Errorf("expected status code %d, got %d", fiber.StatusOK, resp.StatusCode)
				}

			case "exit":
				// Mock the Exit behavior and set a valid cookie
				app.Put("/exit", authCtrl.Exit)

				// Create a new request without body
				req := httptest.NewRequest(tc.method, fmt.Sprintf("/%s", tc.action), nil)

				// Set the session cookie in the request
				cookie := &http.Cookie{
					Name:  "user_session",
					Value: "some_valid_token", // Mocked token value
				}
				req.AddCookie(cookie)

				resp, err := app.Test(req, -1)
				if err != nil {
					t.Fatalf("Fiber test error: %v", err)
				}

				if resp.StatusCode != fiber.StatusOK {
					t.Errorf("expected status code %d, got %d", fiber.StatusOK, resp.StatusCode)
				}
			}
		})
	}
}
