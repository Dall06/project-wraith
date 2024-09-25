package logger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func ReadFile(filePath string) ([]map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	scanner := bufio.NewScanner(file)

	var logs []map[string]interface{} // Slice to hold the parsed log entries
	for scanner.Scan() {
		line := scanner.Text()
		var logEntry map[string]interface{}

		// Unmarshal JSON into map
		err := json.Unmarshal([]byte(line), &logEntry)
		if err != nil {
			log.Printf("Error parsing JSON: %v", err) // Log the error and continue
			continue
		}

		logs = append(logs, logEntry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return logs, nil
}
