package logger

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name            string
		projectPath     string
		encrypt         bool
		encryptKey      string
		setupFileSystem bool
		expectedError   error
	}{
		{
			name:            "Initialize Success",
			projectPath:     "./testlogs",
			setupFileSystem: true,
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFileSystem {
				err := os.MkdirAll(tt.projectPath, os.ModePerm)
				require.NoError(t, err)
				defer func(path string) {
					err := os.RemoveAll(path)
					if err != nil {
						t.Fatalf("Failed to remove temp file: %v", err)
					}
				}(tt.projectPath)
			}

			l := NewLogger(tt.projectPath, tt.encrypt, tt.encryptKey)
			err := l.Initialize()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, l)
			}
		})
	}
}
