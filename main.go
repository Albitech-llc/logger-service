package main

import (
	"log"

	"github.com/Albitech-llc/logger-service/internal/config"
	"github.com/Albitech-llc/logger-service/pkg/caching"
)

func main() {
	// Access configuration statically
	cfg := config.LoadConfig()

	// Initialize Redis
	_, _, err := caching.InitializeRedis(cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	if err != nil {
		log.Fatal("Failed to initialize Redis: ", err)
	}
}
