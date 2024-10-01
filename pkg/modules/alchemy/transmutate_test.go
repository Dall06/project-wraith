package alchemy_test

import (
	"project-wraith/pkg/modules/alchemy"
	"testing"
)

type TestEntity struct {
	Field1 string
	Field2 string
}

func TestTransmutation(t *testing.T) {
	secret := "supersecretkey"
	tests := []struct {
		name     string
		input    TestEntity
		expected TestEntity
	}{
		{
			name: "Test basic encryption",
			input: TestEntity{
				Field1: "value1",
				Field2: "value2",
			},
			expected: TestEntity{}, // We'll decrypt and verify in Revert test
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := alchemy.Transmutation(&tc.input, secret)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			// Just ensure the fields are modified (encrypted)
			if tc.input.Field1 == "value1" || tc.input.Field2 == "value2" {
				t.Errorf("expected fields to be encrypted, but got unmodified values")
			}
		})
	}
}

func TestRevert(t *testing.T) {
	secret := "supersecretkey"

	encryptedValue1, err := alchemy.Encrypt("value1", secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	encryptedValue2, err := alchemy.Encrypt("value2", secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tests := []struct {
		name     string
		input    TestEntity
		expected TestEntity
	}{
		{
			name: "Test basic decryption",
			input: TestEntity{
				Field1: encryptedValue1, // Use actual encrypted values
				Field2: encryptedValue2, // Use actual encrypted values
			},
			expected: TestEntity{
				Field1: "value1",
				Field2: "value2",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := alchemy.Revert(&tc.input, secret)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.input != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, tc.input)
			}
		})
	}
}
