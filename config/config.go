package config

import (
	"log"

	"github.com/joho/godotenv"
)

// Config contains configurable values
type Config struct {
}

var config *Config

const (
	// EnvironmentLocal is local environment for a developers' machine
	EnvironmentLocal = "local"
	// EnvironmentDev is environment for dev server
	EnvironmentDev = "dev"
	// EnvironmentProd is environment for prod server
	EnvironmentProd = "prod"
)

// Init initializes the config with environment (available environment: [development, stage])
func Init(env string) {
	if "" == env {
		env = EnvironmentDev
	}

	log.Printf("Load %v environment", env)

	godotenv.Load(".env." + env)
	godotenv.Load()

	config = &Config{} //TODO: provide api keys and etc.
}

// GetConfig get shared configs
func GetConfig() *Config {
	return config
}
