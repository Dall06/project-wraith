package req

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

func SendRequest(req HTTPRequest) (string, error) {
	client := &http.Client{}

	request, err := http.NewRequest(req.Method, req.URL, bytes.NewBuffer(req.Body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range req.Headers {
		request.Header.Add(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close response body: %w", err)
	}

	return string(body), nil
}
