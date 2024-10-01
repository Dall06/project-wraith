package mail

import (
	"errors"
	"net/smtp"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSMTPClient is a mock type for the smtp client
type MockSMTPClient struct {
	mock.Mock
}

func (m *MockSMTPClient) SendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	args := m.Called(addr, auth, from, to, msg)
	return args.Error(0)
}

func TestSendMail(t *testing.T) {
	cases := []struct {
		name            string
		templateContent string
		templateFile    string
		content         interface{}
		subject         string
		to              []string
		wantError       bool
	}{
		{
			name:            "Template Error",
			templateContent: ``,
			templateFile:    "test_template_error.html",
			content:         nil,
			subject:         "Test Subject",
			to:              []string{"to@example.com"},
			wantError:       true,
		},
		{
			name:            "SMTP Error",
			templateContent: `<!DOCTYPE html><html><body><h1>{{.Title}}</h1></body></html>`,
			templateFile:    "test_template_smtp_error.html",
			content:         map[string]string{"Title": "Hello, World!"},
			subject:         "Test Subject",
			to:              []string{"to@example.com"},
			wantError:       true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSMTP := new(MockSMTPClient)
			mockSMTP.On(
				"SendMail",
				"smtp.example.com:587",
				smtp.PlainAuth("", "from@example.com", "password", "smtp.example.com"),
				"from@example.com", tc.to, mock.Anything).Return(errors.New("smtp error"))

			// Create a mail instance
			m := NewMail("from@example.com", "password", "smtp.example.com", "587")

			// Write the template to a file
			err := writeFile(tc.templateFile, tc.templateContent)
			assert.NoError(t, err)
			defer func(filename string) {
				err := deleteFile(filename)
				if err != nil {
					t.Errorf("Failed to clean up test file: %v", err)
				}
			}(tc.templateFile) // Clean up the file after test

			// Call Send method
			err = m.Send(tc.templateFile, tc.content, tc.subject, tc.to)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Utility functions for file operations
func writeFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func deleteFile(filename string) error {
	return os.Remove(filename)
}
