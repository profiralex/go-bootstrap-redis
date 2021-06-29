package config

import (
	"fmt"

	"github.com/profiralex/goconfig"
)

type Config struct {
	AppConfig struct {
		DebugLevel string `cfg:"DEBUG_LEVEL" cfg-default:"info"`
		Port       string `cfg:"PORT" cfg-default:"8080"`
		ApiKey     string `cfg:"API_KEY" cfg-default:"Please set api key"`
		Version    string `cfg:"APP_VER" cfg-default:"v1"`
	}

	RedisConfig struct {
		Host string `cfg:"REDIS_HOST" cfg-default:"localhost"`
		Port string `cfg:"REDIS_PORT" cfg-default:"6379"`
	}
}

var config *Config = nil

func Init() {
	GetConfig()
}

// GetConfig returns the configuration
func GetConfig() Config {
	if config == nil {
		config = &Config{}
		err := goconfig.Load(config, &goconfig.EnvProvider{}, false)
		if err != nil {
			panic(fmt.Errorf("unable to load configuration: %w", err))
		}
	}
	return *config
}
