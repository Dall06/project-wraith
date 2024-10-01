package link_test

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"project-wraith/pkg/modules/link"
	"testing"
)

func TestError(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name            string
		err             error
		expectedCode    int
		expectedMessage string
	}{
		{
			name:            "Fiber Error",
			err:             fiber.NewError(fiber.StatusNotFound, "not found"),
			expectedCode:    fiber.StatusNotFound,
			expectedMessage: "not found",
		},
		{
			name:            "Generic Error",
			err:             errors.New("an internal error occurred"),
			expectedCode:    fiber.StatusInternalServerError,
			expectedMessage: "an internal error occurred",
		},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := fiber.New()

			// Use the Error function to handle errors
			app.Get("/", func(ctx *fiber.Ctx) error {
				return link.Error(ctx, tt.err)
			})

			req := httptest.NewRequest("GET", "/", nil)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("Fiber test error: %v", err)
			}

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			var response link.Response
			err = json.NewDecoder(resp.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Failed to decode JSON response: %v", err)
			}

			assert.Equal(t, tt.expectedMessage, response.Message)
		})
	}
}
