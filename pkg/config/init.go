package config

import "github.com/spf13/viper"

type Init struct {
	Database struct {
		User struct {
			Uri  string
			Name string
		}
		Manager struct {
			Uri  string
			Name string
		}
		License struct {
			Uri  string
			Name string
		}
	}
	Sms struct {
		ResetAsset string
		From       string
		AccountSID string
		AuthToken  string
	}
	Mail struct {
		From     string
		Password string
		Host     string
		Port     string
	}
	Options struct {
		EncryptResponse bool
		EncryptDbData   bool
		EncryptLogs     bool
		UploadLogs      bool
		UseLicense      bool
	}
}

func LoadInit(fileName, extension, folderPath string) (*Init, error) {
	snake := viper.New()
	snake.SetConfigName(fileName)
	snake.SetConfigType(extension)
	snake.AddConfigPath(folderPath)

	err := snake.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var ini Init
	err = snake.Unmarshal(&ini)
	if err != nil {
		return nil, err
	}

	return &ini, nil
}
