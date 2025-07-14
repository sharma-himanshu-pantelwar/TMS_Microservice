package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	// APP_PORT       string `mapstructure:"APP_PORT"`
	// REDIS_ADDR     string `mapstructure:"REDIS_ADDR"`
	// REDIS_PASSWORD string `mapstructure:"REDIS_PASSWORD"`
	// REDIS_DB       int    `mapstructure:"REDIS_DB"`
}

func LoadConfig() (*Config, error) {
	// Set the path and name for the config file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./.secrets")

	// Set default values
	viper.SetDefault("APP_PORT", "8002")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)

	// Read environment variables
	viper.AutomaticEnv()

	// Try to read the config file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: Could not read config file: %v\n", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return &config, nil
}
