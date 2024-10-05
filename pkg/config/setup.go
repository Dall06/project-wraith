package config

import "github.com/spf13/viper"

type Setup struct {
	Server struct {
		Host               string
		Port               int
		Env                string
		BasePath           string
		CookiesMinutesLife int
	}
	Logger struct {
		Debug      bool
		FolderPath string
	}
	Redirects struct {
		ResetUrl string
	}
}

func LoadSetup(fileName, extension, folderPath string) (*Setup, error) {
	snake := viper.New()
	snake.SetConfigName(fileName)
	snake.SetConfigType(extension)
	snake.AddConfigPath(folderPath)

	if err := snake.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Setup
	if err := snake.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
