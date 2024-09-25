package alchemy_test

import (
	"project-wraith/src/pkg/alchemy"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name     string
		secret   string
		expected int
	}{
		{
			name:     "Test SHA-256 key generation",
			secret:   "supersecretkey",
			expected: 32, // SHA-256 produces a 32-byte key
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			key := alchemy.GenerateKey(tc.secret)
			if len(key) != tc.expected {
				t.Errorf("expected key length %d, got %d", tc.expected, len(key))
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		secret    string
	}{
		{
			name:      "Test basic encryption",
			plaintext: "plaintext",
			secret:    "supersecretkey",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := alchemy.Encrypt(tc.plaintext, tc.secret)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if encrypted == tc.plaintext {
				t.Errorf("expected encrypted value to be different from plaintext")
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	// Encrypt a known plaintext to get the ciphertext
	plaintext := "plaintext"
	secret := "supersecretkey"
	ciphertext, err := alchemy.Encrypt(plaintext, secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tests := []struct {
		name       string
		ciphertext string
		secret     string
		expected   string
	}{
		{
			name:       "Test basic decryption",
			ciphertext: ciphertext, // Use the actual encrypted value
			secret:     secret,
			expected:   plaintext,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			decrypted, err := alchemy.Decrypt(tc.ciphertext, tc.secret)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if decrypted != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, decrypted)
			}
		})
	}
}
