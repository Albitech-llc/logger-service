package logger

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Redis struct {
		Host   string
		Port   string
		LogsDB int
	}

	Logging struct {
		IsPubSub     bool
		LogsFilePath string
		Channels     struct {
			Logs    string
			Info    string
			Warning string
			Error   string
			Debug   string
		}
	}
}

var (
	once   sync.Once
	config *Config
)

func LoadConfig() *Config {
	once.Do(func() {
		env := os.Getenv("ENV")
		if env == "" {
			env = "dev" // Default to development if not set
		}

		fileName := fmt.Sprintf("config.%s.json", env)

		viper.SetConfigName(fileName)
		viper.SetConfigType("json")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("[LOGGER] Error reading config file: %v", err)
		}

		config = &Config{}

		err = viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("[LOGGER] Unable to parse configuration into struct: %v", err)
		}
	})
	return config
}
