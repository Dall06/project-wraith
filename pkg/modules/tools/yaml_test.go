package tools_test

import (
	"io"
	"os"
	"project-wraith/pkg/modules/tools"
	"reflect"
	"testing"
)

func TestReadYaml(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		toBind      interface{}
		expected    interface{}
		expectErr   bool
	}{
		{
			name:        "Valid YAML",
			yamlContent: "key: value\nnumber: 123",
			toBind:      map[string]interface{}{},
			expected:    map[string]interface{}{"key": "value", "number": 123},
			expectErr:   false,
		},
		{
			name:        "Empty File",
			yamlContent: "",
			toBind:      map[string]interface{}{},
			expected:    map[string]interface{}{},
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file with the YAML content
			file, err := os.CreateTemp("", "test*.yaml")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(file.Name())

			// Write the YAML content to the file
			_, err = file.WriteString(tt.yamlContent)
			if err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}

			// Reset the file pointer to the beginning
			file.Seek(0, io.SeekStart)

			// Test ReadYaml function
			err = tools.ReadYaml(file.Name(), tt.toBind)

			if (err != nil) != tt.expectErr {
				t.Errorf("ReadYaml() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			// Check if the binding matches the expected result
			actual := tt.toBind
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("ReadYaml() = %v, want %v", actual, tt.expected)
			}
		})
	}
}
