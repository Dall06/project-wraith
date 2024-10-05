package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func LoadLicense(fileName, extension, folderPath string) (string, error) {
	fileWithExtension := fmt.Sprintf("%s.%s", fileName, extension)
	fullPath := filepath.Join(folderPath, fileWithExtension)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed reading file: %w", err)
	}

	return string(data), nil
}
