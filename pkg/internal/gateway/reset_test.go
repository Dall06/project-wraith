package gateway_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"project-wraith/pkg/internal/gateway"
	"project-wraith/pkg/internal/rules"
	"project-wraith/pkg/modules/logger"
	"project-wraith/pkg/modules/mail"
	"project-wraith/pkg/modules/sms"
	"testing"
)

func TestResetController(test *testing.T) {
	test.Parallel()

	logMock := &logger.MockLogger{}
	resetMock := &rules.MockResetRule{}
	userMock := &rules.MockUserRule{}
	mailMock := &mail.MockMail{}
	smsMock := &sms.MockTwilio{}

	logMock.On("Error", mock.Anything).Return(nil) // Mock error logging
	logMock.On("Info", mock.Anything).Return(nil)  // Mock info logging

	resetCtrl := gateway.NewResetController(
		logMock,
		resetMock,
		userMock,
		"super_secret_key",
		60,
		mailMock,
		smsMock,
		"http://example.com", // Mock web URL
	)

	tests := []struct {
		name       string
		action     string
		method     string
		setupMocks func()
		expectCode int
		request    interface{}
	}{
		{
			name:   "Test Start - Successful",
			action: "start",
			method: "POST",
			request: gateway.Reset{
				Username: "testuser",
				Email:    "test@example.com",
				Phone:    "1234567890",
			},
			expectCode: fiber.StatusAccepted,
			setupMocks: func() {
				resetMock.On("Start", mock.Anything).Return(&rules.Reset{
					Username: "testuser",
					Email:    "test@example.com",
					Phone:    "1234567890",
					Token:    "mock_token",
				}, nil).Once()

				mailMock.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				smsMock.On("SendSMSTwilio", mock.Anything, mock.Anything, mock.Anything).Return("success", nil).Once()
			},
		},
		{
			name:       "Test Start - Failed to Parse Request",
			action:     "start",
			method:     "POST",
			request:    "", // Invalid request body
			expectCode: fiber.StatusBadRequest,
			setupMocks: func() {
				// No need to mock anything as this will fail on parsing the request
			},
		},
		{
			name:   "Test Modify - Successful",
			action: "modify",
			method: "POST",
			request: gateway.Reset{
				NewPassword: "new_secure_password",
			},
			expectCode: fiber.StatusOK,
			setupMocks: func() {
				resetMock.On("Validate", mock.Anything).Return(&rules.Reset{
					ID:    "1",
					Token: "mock_token",
				}, nil).Once()

				userMock.On("Edit", mock.Anything).Return(nil).Once()
			},
		},
		{
			name:   "Test Modify - Invalid Token",
			action: "modify",
			method: "POST",
			request: gateway.Reset{
				NewPassword: "new_secure_password",
			},
			expectCode: fiber.StatusBadRequest,
			setupMocks: func() {
				resetMock.On("Validate", mock.Anything).Return(nil, fmt.Errorf("invalid token")).Once()
			},
		},
	}

	for _, tc := range tests {
		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := fiber.New()
			tc.setupMocks() // Set up the mocks for each test

			switch tc.action {
			case "start":
				app.Post("/reset/start", resetCtrl.Start)
			case "modify":
				app.Post("/reset/modify", resetCtrl.Modify)
			}

			var req *http.Request
			if tc.request != "" {
				body, err := json.Marshal(tc.request)
				if err != nil {
					t.Fatalf("failed to marshal request: %v", err)
				}
				req = httptest.NewRequest(tc.method, fmt.Sprintf("/reset/%s", tc.action), bytes.NewBuffer(body))
			} else {
				req = httptest.NewRequest(tc.method, fmt.Sprintf("/reset/%s", tc.action), nil)
			}

			req.Header.Set("Content-Type", "application/json")
			if tc.action == "modify" {
				req.Header.Set("X-Reset-Token", "mock_token") // Add the reset token for modify action
			}

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("Fiber test error: %v", err)
			}

			assert.Equal(t, tc.expectCode, resp.StatusCode, "Unexpected status code")
		})
	}
}
