package app

import (
	"os"
	"strconv"
	"sync"
)

// DBConfig is the configuration for the DB client.
type DBConfig struct {
	DBConnection string
	DBHost       string
	DBPort       string
	DBDatabase   string
	DBUsername   string
	DBPassword   string
	DBLogLevel   int
}

// AppConfig is the configuration for the DB client.
type AppConfig struct {
	MSTeamsHookURL string
}

// NewAppConfig returns a new Config struct with the configs
func NewAppConfig() *AppConfig {
	return &AppConfig{
		MSTeamsHookURL: os.Getenv("MS_TEAMS_HOOK_URL"),
	}
}

// NewDBConfig returns a new DBConfig.
// DBHosts is a list of DB hosts, separated by commas.
func NewDBConfig() *DBConfig {
	dbLogLevel, _ := strconv.Atoi(os.Getenv("DB_LOG_LEVEL"))
	return &DBConfig{
		DBConnection: os.Getenv("DB_CONNECTION"),
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBDatabase:   os.Getenv("DB_DATABASE"),
		DBUsername:   os.Getenv("DB_USERNAME"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBLogLevel:   dbLogLevel,
	}
}

var (
	configInstance *Config
	configOnce     sync.Once
)

// Config is the configuration struct that returns all the configs
type Config struct {
	DBConfig  *DBConfig
	AppConfig *AppConfig
}

// NewConfig returns a new Config struct with the configs
// Is a singleton with one memory address
func NewConfig() *Config {
	configOnce.Do(func() {
		configInstance = &Config{
			DBConfig:  NewDBConfig(),
			AppConfig: NewAppConfig(),
		}
	})
	return configInstance
}
