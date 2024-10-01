package sms_test

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/url"
	"project-wraith/pkg/modules/req"
	"project-wraith/pkg/modules/sms"
	"testing"
)

// MockHTTPRequester mocks the req.Requester interface
type MockHTTPRequester struct {
	mock.Mock
}

func (m *MockHTTPRequester) SendRequest(req req.HTTPRequest) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

// Test case struct
type testCase struct {
	name             string
	mockResponse     string
	mockError        error
	expectedResponse string
	expectedError    string
}

func TestSendSMSTwilio(t *testing.T) {
	from := "fromNumber"
	accountSID := "accountSID"
	authToken := "authToken"
	asset := "someAsset"
	to := "toNumber"
	args := []string{"arg1", "arg2"}
	expectedMessage := "formattedAssetContent" // Ensure this matches what `tools.FormatAssetContent` generates

	expectedURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSID)
	expectedBody := url.Values{}
	expectedBody.Set("From", from)
	expectedBody.Set("To", to)
	expectedBody.Set("Body", expectedMessage)

	expectedAuth := accountSID + ":" + authToken
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(expectedAuth))

	// Define test cases
	testCases := []testCase{
		{
			name:             "Success Response",
			mockResponse:     "response",
			expectedResponse: "response",
			expectedError:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRequester := new(MockHTTPRequester)

			// Mock SendRequest
			mockRequester.On("SendRequest", req.HTTPRequest{
				Method: "POST",
				URL:    expectedURL,
				Headers: map[string]string{
					"Authorization": "Basic " + encodedAuth,
					"Content-Type":  "application/x-www-form-urlencoded",
				},
				Body: []byte(expectedBody.Encode()),
			}).Return(tc.mockResponse, tc.mockError).Once()

			twilioClient := sms.NewTwilio(from, accountSID, authToken, asset)

			_, err := twilioClient.SendSMSTwilio(to, true, args...)
			if tc.expectedError == "" {
				assert.NoError(t, err)
			}
		})
	}
}
