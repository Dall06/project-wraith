package alchemy_test

import (
	"project-wraith/pkg/modules/alchemy"
	"testing"
)

type TestStruct struct {
	Name  string
	Age   int
	Email string
}

func TestStructIntoString(t *testing.T) {
	secret := "supersecretkey"
	tests := []struct {
		name     string
		input    TestStruct
		expected TestStruct
	}{
		{
			name: "Test Struct to String Conversion",
			input: TestStruct{
				Name:  "Alice",
				Age:   30,
				Email: "alice@example.com",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := alchemy.StructIntoString(tc.input, secret)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if encrypted == "" {
				t.Errorf("expected encrypted string, got empty string")
			}

			// Decrypt and check if the result matches the original struct
			var decryptedStruct TestStruct
			err = alchemy.StringToStruct(encrypted, &secret, &decryptedStruct)
			if err != nil {
				t.Fatalf("expected no error during decryption, got %v", err)
			}

			if decryptedStruct != tc.input {
				t.Errorf("expected struct %v, got %v", tc.input, decryptedStruct)
			}
		})
	}
}

func TestStringToStruct(t *testing.T) {
	secret := "supersecretkey"
	tests := []struct {
		name     string
		input    TestStruct
		secret   *string
		expected TestStruct
	}{
		{
			name: "Test String to Struct Conversion",
			input: TestStruct{
				Name:  "Bob",
				Age:   25,
				Email: "bob@example.com",
			},
			secret: &secret,
			expected: TestStruct{
				Name:  "Bob",
				Age:   25,
				Email: "bob@example.com",
			},
		},
		{
			name: "Test Nil Secret",
			input: TestStruct{
				Name:  "Charlie",
				Age:   40,
				Email: "charlie@example.com",
			},
			secret: nil,
			expected: TestStruct{
				Name:  "Charlie",
				Age:   40,
				Email: "charlie@example.com",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Encrypt the input struct only if the secret is not nil
			var encrypted string
			var err error
			if tc.secret != nil {
				encrypted, err = alchemy.StructIntoString(tc.input, *tc.secret)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}

			// Test the StringToStruct function
			var output TestStruct
			err = alchemy.StringToStruct(encrypted, tc.secret, &output)

			if tc.secret == nil {
				if err == nil || err.Error() != "secret cannot be nil" {
					t.Errorf("expected error 'secret cannot be nil', got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if output != tc.expected {
				t.Errorf("expected struct %v, got %v", tc.expected, output)
			}
		})
	}
}
