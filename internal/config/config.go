package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	RedisHost string
	RedisPort string
	RedisDB   int
}

var (
	once   sync.Once
	config *Config
)

func LoadConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		config = &Config{
			RedisHost: viper.GetString("redis_host"),
			RedisPort: viper.GetString("redis_port"),
			RedisDB:   viper.GetInt("redis_db"),
		}
	})
	return config
}
