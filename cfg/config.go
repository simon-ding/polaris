package cfg

import (
	"polaris/log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	TMDB TMDB `mapstructure:"tmdb"`
}

type TMDB struct {
	ApiKey string `mapstructure:"apiKey"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app/data")

	var cc Config
	// optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("create config file")
			viper.SafeWriteConfig()
		} else {
			// Config file was found but another error was produced
		}

		return nil, errors.Wrap(err, "load config")
	}

	if err := viper.Unmarshal(&cc); err != nil {
		return nil, errors.Wrap(err, "unmarshal file")
	}
	return &cc, err
}
