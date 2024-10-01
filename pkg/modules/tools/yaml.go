package tools

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

func ReadYaml(filePath string, toBind interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, toBind)
	if err != nil {
		return err
	}

	return nil
}
