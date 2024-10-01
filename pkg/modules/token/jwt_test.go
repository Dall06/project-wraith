package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

func TestCreateJwtToken(t *testing.T) {

	testCases := []struct {
		secret string
		exp    time.Duration
		data   interface{}
	}{
		{"secret1", time.Minute, "testdata1"},
		{"secret2", time.Hour, map[string]interface{}{"key": "value"}},
		{"secret3", time.Second, 12345},
	}

	for _, tc := range testCases {
		token, err := CreateJwtToken(tc.secret, tc.exp, tc.data)
		if err != nil {
			t.Errorf("failed to create JWT token: %v", err)
		}

		if token == "" {
			t.Error("expected non-empty token")
		}
	}
}

func TestExpireJwtToken(t *testing.T) {
	testCases := []struct {
		secret string
		exp    time.Duration
		data   interface{}
	}{
		{"secret1", time.Minute, "testdata1"},
		{"secret2", time.Hour, map[string]interface{}{"key": "value"}},
		{"secret3", time.Second, 12345},
	}

	for _, tc := range testCases {
		token, err := ExpireJwtToken(tc.secret, tc.exp, tc.data)
		if err != nil {
			t.Errorf("failed to expire JWT token: %v", err)
		}

		if token == "" {
			t.Error("expected non-empty token")
		}
	}
}

func TestValidateJwtToken(t *testing.T) {
	testCases := []struct {
		secret          string
		exp             time.Duration
		data            interface{}
		extraValidation func(jwt.MapClaims) error
		expectValid     bool
	}{
		{"secret1", time.Minute, "testdata1", nil, true},
		{"secret2", time.Second, "testdata2", nil, false}, // Should expire quickly
		{"secret3", time.Hour, map[string]interface{}{"key": "value"}, func(claims jwt.MapClaims) error {
			if claims["data"] != "testdata3" {
				return errors.New("invalid data")
			}
			return nil
		}, false},
	}

	for _, tc := range testCases {
		token, err := CreateJwtToken(tc.secret, tc.exp, tc.data)
		if err != nil {
			t.Fatalf("failed to create JWT token: %v", err)
		}

		time.Sleep(2 * time.Second) // Sleep to allow short-lived tokens to expire

		valid, err := ValidateJwtToken(token, tc.secret, tc.extraValidation)
		if valid != tc.expectValid {
			t.Errorf("expected valid: %v, got: %v, err: %v", tc.expectValid, valid, err)
		}
	}
}
