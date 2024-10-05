package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Secrets struct {
	Server struct {
		KeyWord string `mapstructure:"SERVER_KEY_WORD"`
	}
	Secrets struct {
		Jwt       string `mapstructure:"SECRET_JWT"`
		DbData    string `mapstructure:"SECRET_DB"`
		Response  string `mapstructure:"SECRET_RESPONSE"`
		Password  string `mapstructure:"SECRET_PASSWORD"`
		Cookies   string `mapstructure:"SECRET_COOKIES"`
		Internals string `mapstructure:"SECRET_INTERNALS"`
		Logs      string `mapstructure:"SECRET_LOGS"`
	}
	Notifiers struct {
		Bot struct {
			Token string `mapstructure:"NOTIFIER_TLG_BOT_TOKEN"`
			Chat  string `mapstructure:"NOTIFIER_TLG_BOT_CHAT"`
		}
	}
}

func LoadSecrets(fileName, extension, folderPath string) (*Secrets, error) {
	snake := viper.New()
	snake.SetConfigName(fileName)
	snake.SetConfigType(extension)
	snake.AddConfigPath(folderPath)

	err := snake.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	snake.AutomaticEnv()

	var secrets Secrets
	err = snake.Unmarshal(&secrets)
	if err != nil {
		return nil, err
	}

	return &secrets, nil
}
