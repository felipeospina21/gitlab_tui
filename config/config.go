package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type config struct {
	BaseURL    string `mapstructure:"base_url"`
	APIToken   string `mapstructure:"token"`
	APIVersion string `mapstructure:"api_version"`
}

var Config config

func Load(configObj *config) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AddConfigPath("$HOME/.config/glabt/")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshal: %w", err))
	}
}
