package tools

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadAsset(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	var contentBuilder strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		contentBuilder.WriteString(scanner.Text())
		contentBuilder.WriteString("\n")
	}

	errClose := file.Close()
	if errClose != nil {
		return "", errClose
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return contentBuilder.String(), nil
}

func FormatAssetContent(content string, args ...interface{}) string {
	return fmt.Sprintf(content, args...)
}
