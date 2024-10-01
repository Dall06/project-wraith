package tools_test

import (
	"os"
	"project-wraith/pkg/modules/tools"
	"testing"
)

func TestReadAsset(t *testing.T) {
	// Define the test cases using a slice of structs
	tests := []struct {
		name        string
		fileContent string
		shouldError bool
	}{
		{
			name:        "ValidFile",
			fileContent: "Hello, World!\nThis is a test file.",
			shouldError: false,
		},
		{
			name:        "NonExistentFile",
			fileContent: "",
			shouldError: true,
		},
	}

	// Iterate over test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var tempFileName string
			if !tc.shouldError {
				// Create a temporary file with test content
				tempFile, err := os.CreateTemp("", "test_asset_*.txt")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				defer func(name string) {
					err := os.Remove(name)
					if err != nil {
						t.Fatalf("Failed to remove temp file: %v", err)
					}
				}(tempFile.Name())
				tempFileName = tempFile.Name()

				if _, err := tempFile.WriteString(tc.fileContent); err != nil {
					t.Fatalf("Failed to write to temp file: %v", err)
				}

				// Close the temp file so it can be opened by ReadAsset
				err = tempFile.Close()
				if err != nil {
					t.Fatalf("Failed to close temp file: %v", err)
				}
			} else {
				// Non-existent file case
				tempFileName = "non_existent_file.txt"
			}

			// Test ReadAsset function
			readContent, err := tools.ReadAsset(tempFileName)

			if tc.shouldError {
				if err == nil {
					t.Errorf("Expected an error for test case %s, but got nil", tc.name)
				}
			} else {
				if err != nil {
					t.Fatalf("ReadAsset returned an error: %v", err)
				}

				// Check if the content matches
				if readContent != tc.fileContent+"\n" { // \n is added by scanner.Text()
					t.Errorf("Expected content: %v, but got: %v", tc.fileContent, readContent)
				}
			}
		})
	}
}

func TestFormatAssetContent(t *testing.T) {
	// Define the test cases using a slice of structs
	tests := []struct {
		name     string
		content  string
		args     []interface{}
		expected string
	}{
		{
			name:     "BasicFormatting",
			content:  "Hello, %s! Today is %s.",
			args:     []interface{}{"Alice", "Monday"},
			expected: "Hello, Alice! Today is Monday.",
		},
		{
			name:     "NumberFormatting",
			content:  "Number: %d, Float: %.2f",
			args:     []interface{}{42, 3.14159},
			expected: "Number: 42, Float: 3.14",
		},
	}

	// Iterate over test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formatted := tools.FormatAssetContent(tc.content, tc.args...)
			if formatted != tc.expected {
				t.Errorf("Expected: %v, but got: %v", tc.expected, formatted)
			}
		})
	}
}
