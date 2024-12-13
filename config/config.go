package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Port       string `envconfig:"port" default:"8080"`
	Env        string `envconfig:"env" default:"local"`
	DbHost     string `envconfig:"db_host" default:"localhost"`
	DbPort     int    `envconfig:"db_port" default:"5432"`
	DbUsername string `envconfig:"db_username" default:"spankie"`
	DbPassword string `envconfig:"db_password" default:"spankie"`
	DbName     string `envconfig:"db_name" default:"spankie"`
	LogLevel   string `envconfig:"log_level" default:"info"`
}

// GetConfig create a configuration object from the environment variables,
// uses the configuration to set up the logger and returns the configuration object
func GetConfig() (*Configuration, error) {
	config := &Configuration{}
	err := envconfig.Process("gymshark", config)
	if err != nil {
		return nil, err
	}

	// override port if set in env
	envPort := os.Getenv("PORT")
	if envPort != "" {
		config.Port = envPort
	}

	return config, nil
}
