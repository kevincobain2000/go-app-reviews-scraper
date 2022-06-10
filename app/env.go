package app

import (
	"os"

	"github.com/joho/godotenv"
)

// SetEnv will set envs paths and load .env file
// If ENV_PATH variable is passed, it will be used as env path
// ENV_PATH=/full/path/.env main 3000
func SetEnv() {
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = ".env"
	}
	err := godotenv.Load(envPath)
	if err != nil {
		panic("fatal: error loading .env file")
	}
}
