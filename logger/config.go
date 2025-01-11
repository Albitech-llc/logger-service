package logger

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	IsPubSub     bool
	RedisHost    string
	RedisPort    string
	RedisDB      int
	LogsFilePath string

	LogsChannel    string
	InfoChannel    string
	WarningChannel string
	ErrorChannel   string
	DebugChannel   string
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
			IsPubSub:     viper.GetBool("is_pub_sub"),
			RedisHost:    viper.GetString("redis_host"),
			RedisPort:    viper.GetString("redis_port"),
			RedisDB:      viper.GetInt("redis_db"),
			LogsFilePath: viper.GetString("logs_file_path"),

			LogsChannel:    viper.GetString("logs_channel"),
			InfoChannel:    viper.GetString("info_channel"),
			WarningChannel: viper.GetString("warning_channel"),
			ErrorChannel:   viper.GetString("error_channel"),
			DebugChannel:   viper.GetString("debug_channel"),
		}
	})
	return config
}
