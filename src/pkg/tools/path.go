package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetProjectPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	for dir := cwd; dir != string(filepath.Separator); dir = filepath.Dir(dir) {
		_, err := os.Stat(filepath.Join(dir, "go.mod"))
		if err == nil {
			return dir, nil
		}

		if !os.IsNotExist(err) {
			return "", fmt.Errorf("failed to check directory: %w", err)
		}
	}

	return "", fmt.Errorf("failed to find project root directory")
}

func BuildPath(name, extension string, folderPath ...string) string {
	path := filepath.Join(append([]string{name}, folderPath...)...)
	return fmt.Sprintf("%s.%s", path, extension)
}
