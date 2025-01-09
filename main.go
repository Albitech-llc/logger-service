package main

import (
	"log"

	"github.com/Albitech-llc/logger-service/logger"
	"github.com/Albitech-llc/logger-service/pkg/caching"
)

func main() {
	// Access configuration statically
	cfg := logger.LoadConfig()

	// Initialize Redis
	_, _, err := caching.InitializeRedis(cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	if err != nil {
		log.Printf("failed to initialize Redis: %v\n", err)
	}
}
