package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

type Secrets struct {
	Server struct {
		KeyWord string
	}
	Keys struct {
		Jwt       string
		DbData    string
		Response  string
		Password  string
		Cookies   string
		Internals string
		Logs      string
	}
	Storage struct {
		AccessKey string
		SecretKey string
	}
	Notifiers struct {
		Bot struct {
			Token string
			Chat  string
		}
	}
}

func LoadSecrets(fileName, extension, folderPath string) (*Secrets, error) {
	fileNameWithExt := fmt.Sprintf("%s.%s", fileName, extension)
	filePath := filepath.Join(folderPath, fileNameWithExt)

	err := godotenv.Load(filePath)
	if err != nil {
		fmt.Printf("could not load env file, falling back to system environment variables... %v\n", err)
	}

	var secrets Secrets

	secrets.Server.KeyWord = os.Getenv("SERVER_KEY_WORD")
	secrets.Keys.Jwt = os.Getenv("SECRET_JWT")
	secrets.Keys.DbData = os.Getenv("SECRET_DB")
	secrets.Keys.Response = os.Getenv("SECRET_RESPONSE")
	secrets.Keys.Password = os.Getenv("SECRET_PASSWORD")
	secrets.Keys.Cookies = os.Getenv("SECRET_COOKIES")
	secrets.Keys.Internals = os.Getenv("SECRET_INTERNALS")
	secrets.Keys.Logs = os.Getenv("SECRET_LOGS")
	secrets.Notifiers.Bot.Token = os.Getenv("NOTIFIER_TLG_BOT_TOKEN")
	secrets.Notifiers.Bot.Chat = os.Getenv("NOTIFIER_TLG_BOT_CHAT")

	return &secrets, nil
}
