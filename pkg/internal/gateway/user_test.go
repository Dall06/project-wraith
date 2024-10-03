package gateway_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project-wraith/pkg/internal/gateway"
	"project-wraith/pkg/internal/rules"
	"project-wraith/pkg/modules/logger"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserController(test *testing.T) {
	test.Parallel()

	logMock := &logger.MockLogger{}
	userMock := &rules.MockUserRule{}

	logMock.On("Error", mock.Anything).Return(nil)
	logMock.On("Info", mock.Anything).Return(nil)
	logMock.On("Warn", mock.Anything).Return(nil)

	controller := gateway.NewUserController(logMock, userMock, "secret", false, "responseSecret", 60)

	tests := []struct {
		name             string
		method           string
		url              string
		request          interface{}
		expectedStatus   int
		expectedResponse map[string]interface{}
		setupMocks       func()
	}{
		{
			name:   "Test Register - Successful",
			method: "POST",
			url:    "/user/register",
			request: rules.User{
				Username: "newuser",
				Password: "newpassword",
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "register successful",
			},
			setupMocks: func() {
				userMock.On("Register", mock.Anything).Return(&rules.User{
					Username: "newuser",
				}, nil).Once()
			},
		},
		{
			name:           "Test Get User Details - Successful",
			method:         "GET",
			url:            "/user/newuser",
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"content": map[string]interface{}{
					"email": "",
					"id":    "newuser",
					"name":  "", "phone": "",
					"username": "newuser",
				}},
			setupMocks: func() {
				userMock.On("Get", mock.Anything).Return(&rules.User{
					ID:       "newuser",
					Username: "newuser",
				}, nil).Once()
			},
		},
		{
			name:   "Test Edit User - Successful",
			method: "PUT",
			url:    "/user/edit",
			request: rules.User{
				Username: "updateduser",
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "edit successful",
			},
			setupMocks: func() {
				userMock.On("Edit", mock.Anything).Return(nil).Once()
			},
		},
		{
			name:   "Test Remove User - Successful",
			method: "DELETE",
			url:    "/user/remove",
			request: rules.User{
				ID: "testuser",
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "remove successful",
			},
			setupMocks: func() {
				userMock.On("Remove", mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, testCase := range tests {
		test.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			app := fiber.New()
			testCase.setupMocks() // Set up the mocks for each test

			switch testCase.method {
			case "POST":
				if testCase.url == "/user/register" {
					app.Post(testCase.url, controller.Register)
				}
			case "GET":
				app.Get("/user/:id", controller.Get)
			case "PUT":
				app.Put(testCase.url, controller.Edit)
			case "DELETE":
				app.Delete(testCase.url, controller.Remove)
			}

			// Create the request payload
			var req *http.Request
			if testCase.request != nil {
				payload, _ := json.Marshal(testCase.request)
				req = httptest.NewRequest(testCase.method, testCase.url, bytes.NewBuffer(payload))
			} else {
				req = httptest.NewRequest(testCase.method, testCase.url, nil)
			}

			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}

			// Validate response
			assert.Equal(t, testCase.expectedStatus, resp.StatusCode)
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, testCase.expectedResponse, response)
		})
	}
}
