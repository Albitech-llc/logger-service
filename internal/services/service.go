// Package services contains the business logic implementations.
// It interacts with repository layers to perform operations on data.
package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Albitech-llc/logger-service/internal/config"
	"github.com/Albitech-llc/logger-service/internal/models"              // Importing model
	"github.com/Albitech-llc/logger-service/internal/services/interfaces" // Importing service interface definitions
	"github.com/Albitech-llc/logger-service/pkg/caching"
)

// Service is a struct that holds a reference to the RepositoryI interface,
// which is used to interact with the data in the database.
type service struct {
}

// NewService initializes and returns an instance of service,
// implementing the ServiceI interface.
func NewService() interfaces.Service {
	return &service{} // Dependency injection of repository
}

func (s *service) LogInfo(service string, message string) error {
	log := models.LogMessage{
		Level:     "INFO",
		Timestamp: time.Now(),
		Service:   service,
		Message:   message,
	}

	//Serialize log object to JSON
	logJSON, err := json.Marshal(log)
	if err != nil {
		fmt.Printf("Error serializing log: %v", err)
		return fmt.Errorf("error serializing log")
	}

	err = publish(logJSON, "logs")
	if err != nil {
		fmt.Printf("Failed to publish: %v", err)
	}
	publish(logJSON, "log-info")

	return nil
}

func (s *service) LogWarning(service string, message string) error {
	log := models.LogMessage{
		Level:     "WARN",
		Timestamp: time.Now(),
		Service:   service,
		Message:   message,
	}

	//Serialize log object to JSON
	logJSON, err := json.Marshal(log)
	if err != nil {
		fmt.Printf("Error serializing log: %v", err)
		return fmt.Errorf("error serializing log")
	}

	err = publish(logJSON, "logs")
	if err != nil {
		fmt.Printf("Failed to publish: %v", err)
	}
	publish(logJSON, "log-warn")

	return nil
}

func (s *service) LogError(service string, message string) error {
	log := models.LogMessage{
		Level:     "ERROR",
		Timestamp: time.Now(),
		Service:   service,
		Message:   message,
	}

	//Serialize log object to JSON
	logJSON, err := json.Marshal(log)
	if err != nil {
		fmt.Printf("Error serializing log: %v", err)
		return fmt.Errorf("error serializing log")
	}

	err = publish(logJSON, "logs")
	if err != nil {
		fmt.Printf("Failed to publish: %v", err)
	}
	publish(logJSON, "log-error")

	return nil
}

func (s *service) LogMessage(service string, message string, level string) error {
	log := models.LogMessage{
		Level:     level,
		Timestamp: time.Now(),
		Service:   service,
		Message:   message,
	}

	//Serialize log object to JSON
	logJSON, err := json.Marshal(log)
	if err != nil {
		fmt.Printf("Error serializing log: %v", err)
		return fmt.Errorf("error serializing log")
	}

	err = publish(logJSON, "logs")
	if err != nil {
		fmt.Printf("Failed to publish: %v", err)
	}
	publish(logJSON, "log-custom")

	return nil
}

func publish(json []byte, channel string) error {
	// Access configuration statically
	cfg := config.LoadConfig()

	// Initialize Redis
	rdb, ctx, err := caching.InitializeRedis(cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	if err != nil {
		fmt.Printf("Failed to initialize Redis: %v", err)
		return fmt.Errorf("failed to initialize Redis")
	}

	// Use a goroutine for asynchronous publishing
	go func() {
		err = rdb.Publish(ctx, channel, json).Err()
		if err != nil {
			fmt.Printf("Failed to publish to Redis channel '%s': %v\n", channel, err)
		}
	}()

	return nil
}
