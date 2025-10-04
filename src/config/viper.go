package config

import (
	"log"

	"fmt"

	"github.com/joho/godotenv"

	"github.com/spf13/viper"
)

// NewViper is a function to load config from .env file
func NewViper() *viper.Viper {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := viper.New()

	// Set config to read from environment variables
	config.AutomaticEnv()
	fmt.Println("Loading configuration from .env file...")
	config.SetConfigFile(".env")

	// Read the config file
	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	return config
}
