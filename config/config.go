package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config contains configurable values
type Config struct {
	DatabaseHost string
	DatabasePort string
	DatabaseName string
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

	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err.Error())
	}
	godotenv.Load()

	config = &Config{
		DatabaseHost: os.Getenv("DATABASE_HOST"),
		DatabasePort: os.Getenv("DATABASE_PORT"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
	}

	log.Printf("%v", config.DatabasePort)
}

// GetConfig get shared configs
func GetConfig() *Config {
	return config
}
