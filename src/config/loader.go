package config

import "github.com/spf13/viper"

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
