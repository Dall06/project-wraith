package req_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"project-wraith/pkg/modules/req"
	"testing"
)

func TestSendRequest(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name             string
		mockResponse     string
		mockStatusCode   int
		mockHeaders      map[string]string
		request          req.HTTPRequest
		expectedResponse string
		expectedError    string
	}{
		{
			name:             "Successful Request",
			mockResponse:     "success",
			mockStatusCode:   http.StatusOK,
			mockHeaders:      map[string]string{"Content-Type": "text/plain"},
			request:          req.HTTPRequest{Method: "GET", URL: "", Headers: map[string]string{"Content-Type": "text/plain"}, Body: nil},
			expectedResponse: "success",
			expectedError:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock server that will return the predefined response
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tc.request.Method {
					t.Fatalf("expected method %s, got %s", tc.request.Method, r.Method)
				}
				if r.URL.Path != "/" {
					t.Fatalf("expected URL path /, got %s", r.URL.Path)
				}
				for key, value := range tc.mockHeaders {
					if r.Header.Get(key) != value {
						t.Fatalf("expected header %s with value %s, got %s", key, value, r.Header.Get(key))
					}
				}
				w.WriteHeader(tc.mockStatusCode)
				w.Write([]byte(tc.mockResponse))
			}))
			defer mockServer.Close()

			// Update the URL of the request to use the mock server URL
			tc.request.URL = mockServer.URL

			// Call SendRequest
			response, err := req.SendRequest(tc.request)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, response)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			}
		})
	}
}
