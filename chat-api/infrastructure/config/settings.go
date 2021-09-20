package config

import (
	"os"
)

// Config represents the configuration parameters for the app
type Config struct {
	AccessSecret string
}

// GetConfig returns Configuration items
func GetConfig() *Config {

	return &Config{
		AccessSecret: Getenv("ACCESS_SECRET", ""),
	}
}

// Getenv gets particular env value
func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
