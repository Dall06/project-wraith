package tools_test

import (
	"os"
	"project-wraith/src/pkg/tools"
	"strings"
	"testing"
)

func TestGetProjectPath(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() string
		expectedError bool
		expectedPart  string
	}{
		{
			name: "ValidProjectPath",
			setup: func() string {
				// Set the current directory to a known Go project directory
				originalDir, _ := os.Getwd()
				return originalDir
			},
			expectedError: false,
			// We expect the project path to contain this folder (e.g., "project-wraith")
			expectedPart: "project-wraith",
		},
		{
			name: "InvalidProjectPath",
			setup: func() string {
				// Simulate a directory with no go.mod file
				tempDir := os.TempDir()
				_ = os.Chdir(tempDir)
				return tempDir
			},
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			originalDir := tc.setup()
			defer func(dir string) {
				err := os.Chdir(dir)
				if err != nil {
					t.Errorf("Failed to restore original directory: %v", err)
				}
			}(originalDir) // Restore the original directory after test

			projectPath, err := tools.GetProjectPath()

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// Check that the project path contains the expected part
				if !strings.Contains(projectPath, tc.expectedPart) {
					t.Errorf("Expected project path to contain: %v, but got: %v", tc.expectedPart, projectPath)
				}
			}
		})
	}
}

func TestBuildPath(t *testing.T) {
	tests := []struct {
		name       string
		fileName   string
		extension  string
		folderPath string
		expected   string
	}{
		{
			name:       "SimplePath",
			fileName:   "file",
			extension:  "txt",
			folderPath: "",
			expected:   "file.txt",
		},
		{
			name:       "NestedPath",
			fileName:   "file",
			extension:  "txt",
			folderPath: "folder1/folder2",
			expected:   "folder1/folder2/file.txt",
		},
		{
			name:       "PathWithExtension",
			fileName:   "file.name",
			extension:  "txt",
			folderPath: "folder",
			expected:   "folder/file.name.txt",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tools.BuildPath(tc.fileName, tc.extension, tc.folderPath)
			if result != tc.expected {
				t.Errorf("Expected: %v, but got: %v", tc.expected, result)
			}
		})
	}
}
