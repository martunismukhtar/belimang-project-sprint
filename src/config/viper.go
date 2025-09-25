package config

import (
	"github.com/spf13/viper"
)

// NewViper is a function to load config from .env file
func NewViper() *viper.Viper {
	config := viper.New()

	// Set config to read from .env file
	config.AutomaticEnv()

	return config
}
