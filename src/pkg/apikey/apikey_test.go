package apikey

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestCrateApiKey(t *testing.T) {
	type testCase struct {
		secret   string
		expected string
	}

	testCases := []testCase{
		{"secret1", hex.EncodeToString(sha256.New().Sum([]byte("secret1")))},
		{"secret2", hex.EncodeToString(sha256.New().Sum([]byte("secret2")))},
		{"longersecret", hex.EncodeToString(sha256.New().Sum([]byte("longersecret")))},
	}

	for _, tc := range testCases {
		result := CrateApiKey(tc.secret)
		expectedHash := sha256.New()
		expectedHash.Write([]byte(tc.secret))
		expected := hex.EncodeToString(expectedHash.Sum(nil))

		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	}
}
