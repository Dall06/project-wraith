package config

import "github.com/spf13/viper"

type Config struct {
	Server struct {
		Host               string
		Port               int
		Name               string
		Header             string
		Env                string
		BasePath           string
		KeyWord            string
		CookiesMinutesLife int
		License            string
	}
	Database struct {
		Uri  string
		Name string
	}
	Logger struct {
		Debug      bool
		FolderPath string
		Upload     bool
	}
	Options struct {
		EncryptResponse bool
		EncryptDbData   bool
	}
	Redirects struct {
		ResetUrl string
	}
}

func Load(fileName, extension, folderPath string) (*Config, error) {
	viper.SetConfigName(fileName)
	viper.SetConfigType(extension)
	viper.AddConfigPath(folderPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
