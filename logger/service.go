// Package services contains the business logic implementations.
// It interacts with repository layers to perform operations on data.
package logger

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/Albitech-llc/logger-service/logger/helper"
	"github.com/Albitech-llc/logger-service/logger/models"
	"github.com/Albitech-llc/logger-service/pkg/caching"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// Service is a struct that holds a reference to the RepositoryI interface,
// which is used to interact with the data in the database.
type service struct {
	logger       *logrus.Logger
	logChan      chan []byte
	logInfoChan  chan []byte
	logWarnChan  chan []byte
	logErrorChan chan []byte
	logDebugChan chan []byte
}

// NewService initializes and returns an instance of service,
func NewService() *service {
	// Initialize a logrus logger
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}) // Use JSON formatter for structured logging
	log.SetLevel(logrus.InfoLevel) // Set default log level

	s := &service{
		logger:       log,
		logChan:      make(chan []byte, 10000), // Buffered channel for logs
		logInfoChan:  make(chan []byte, 8000),  // Buffered channel for info logs
		logWarnChan:  make(chan []byte, 1000),  // Buffered channel for warning logs
		logErrorChan: make(chan []byte, 500),   // Buffered channel for error logs
		logDebugChan: make(chan []byte, 500),   // Buffered channel for debug/custom logs
	}
	go s.startPublish()
	return s
}

func (s *service) LogInfo(service string, message string) {
	log := models.LogMessage{
		Level:     LevelInfo,
		Timestamp: time.Now(),
		Service:   service,
		Message:   message,
	}

	s.queueInfoLog(log)
	s.queueLog(log)
}

func (s *service) LogWarning(service string, message string) {
	log := models.LogMessage{
		Level:     LevelWarning,
		Timestamp: time.Now(),
		Service:   service,
		Message:   message,
	}

	s.queueWarningLog(log)
	s.queueLog(log)
}

func (s *service) LogError(service string, message string) {
	log := models.LogMessage{
		Level:     LevelError,
		Timestamp: time.Now(),
		Service:   service,
		Message:   message,
	}

	s.queueErrorLog(log)
	s.queueLog(log)
}

func (s *service) LogMessage(service string, message string, level string) {
	log := models.LogMessage{
		Level:     level,
		Timestamp: time.Now(),
		Service:   service,
		Message:   message,
	}

	s.queueDebugLog(log)
	s.queueLog(log)
}

func (s *service) queueLog(log models.LogMessage) {
	//Serialize log object to JSON
	logJSON, err := json.Marshal(log)
	if err != nil {
		s.logger.SetLevel(logrus.ErrorLevel)
		s.logger.WithError(err).Error("Error serializing log")
		return
	}

	select {
	case s.logChan <- logJSON:
		// Log successfully queued
		return
	default:
		s.logger.SetLevel(logrus.WarnLevel)
		s.logger.Warn("Logs channel is full, dropping log")
		// Handle case where channel is full (e.g., drop or retry)
		time.Sleep(10 * time.Millisecond) // Retry after a short delay
	}
}

func (s *service) queueInfoLog(log models.LogMessage) {
	//Serialize log object to JSON
	logJSON, err := json.Marshal(log)
	if err != nil {
		s.logger.SetLevel(logrus.ErrorLevel)
		s.logger.WithError(err).Error("Error serializing log")
		return
	}

	select {
	case s.logInfoChan <- logJSON:
		// Log successfully queued
		return
	default:
		s.logger.SetLevel(logrus.WarnLevel)
		s.logger.Warn("INFO Log channel is full, dropping log")
		// Handle case where channel is full (e.g., drop or retry)
		time.Sleep(10 * time.Millisecond) // Retry after a short delay
	}
}

func (s *service) queueWarningLog(log models.LogMessage) {
	//Serialize log object to JSON
	logJSON, err := json.Marshal(log)
	if err != nil {
		s.logger.SetLevel(logrus.ErrorLevel)
		s.logger.WithError(err).Error("Error serializing log")
		return
	}

	select {
	case s.logWarnChan <- logJSON:
		// Log successfully queued
		return
	default:
		s.logger.SetLevel(logrus.WarnLevel)
		s.logger.Warn("WARNING Log channel is full, dropping log")
		// Handle case where channel is full (e.g., drop or retry)
		time.Sleep(10 * time.Millisecond) // Retry after a short delay
	}
}

func (s *service) queueErrorLog(log models.LogMessage) {
	//Serialize log object to JSON
	logJSON, err := json.Marshal(log)
	if err != nil {
		s.logger.SetLevel(logrus.ErrorLevel)
		s.logger.WithError(err).Error("Error serializing log")
		return
	}

	select {
	case s.logErrorChan <- logJSON:
		// Log successfully queued
		return
	default:
		s.logger.SetLevel(logrus.WarnLevel)
		s.logger.Warn("ERROR Log channel is full, dropping log")
		// Handle case where channel is full (e.g., drop or retry)
		time.Sleep(10 * time.Millisecond) // Retry after a short delay
	}
}

func (s *service) queueDebugLog(log models.LogMessage) {
	//Serialize log object to JSON
	logJSON, err := json.Marshal(log)
	if err != nil {
		s.logger.SetLevel(logrus.ErrorLevel)
		s.logger.WithError(err).Error("Error serializing log")
		return
	}

	select {
	case s.logDebugChan <- logJSON:
		// Log successfully queued
		return
	default:
		s.logger.SetLevel(logrus.WarnLevel)
		s.logger.Warn("DEBUG Log channel is full, dropping log")
		// Handle case where channel is full (e.g., drop or retry)
		time.Sleep(10 * time.Millisecond) // Retry after a short delay
	}
}

func (s *service) startPublish() {
	// Access configuration statically
	cfg := LoadConfig()

	var rdb *redis.Client = nil
	var ctx context.Context
	var err error

	if cfg.IsPubSub {
		// Initialize Redis
		rdb, ctx, err = caching.InitializeRedis(cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
		if err != nil {
			s.logger.WithError(err).Error("Failed to initialize Redis. Falling back to file logging")
			rdb = nil // Mark Redis as unavailable
		}

		if rdb != nil {
			defer rdb.Close()
		}
	}

	// Create a fallback log file
	fallbackFile, fileErr := os.OpenFile(cfg.LogsFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		s.logger.WithError(fileErr).Error("Failed to create fallback log file")
	}
	//defer fallbackFile.Close()

	go func() {
		for log := range s.logChan {
			s.logger.Infoln(string(log))
			if rdb != nil {
				err := rdb.Publish(ctx, "logs", log).Err()
				if err != nil {
					s.logger.WithError(err).Error("Failed to publish log to Redis. Writing to file")
					helper.WriteToFile(fallbackFile, log) // Fallback to file logging
				}
			} else {
				// Redis unavailable, directly write to file
				helper.WriteToFile(fallbackFile, log)
			}
		}
	}()

	go func() {
		for log := range s.logInfoChan {
			s.logger.Infof("\033[34m%s\033[0m \n", string(log))

			if rdb != nil {
				err := rdb.Publish(ctx, cfg.InfoChannel, log).Err()
				if err != nil {
					s.logger.WithError(err).Error("Failed to publish log to Redis. Writing to file")
					helper.WriteToFile(fallbackFile, log) // Fallback to file logging
				}
			} else {
				// Redis unavailable, directly write to file
				helper.WriteToFile(fallbackFile, log)
			}
		}
	}()

	go func() {
		for log := range s.logWarnChan {
			s.logger.Warningf("\033[33m%s\033[0m \n", string(log))

			if rdb != nil {
				err := rdb.Publish(ctx, cfg.WarningChannel, log).Err()
				if err != nil {
					s.logger.WithError(err).Error("Failed to publish log to Redis. Writing to file")
					helper.WriteToFile(fallbackFile, log) // Fallback to file logging
				}
			} else {
				// Redis unavailable, directly write to file
				helper.WriteToFile(fallbackFile, log)
			}
		}
	}()

	go func() {
		for log := range s.logErrorChan {
			s.logger.Errorf("\033[31m%s\033[0m \n", string(log))

			if rdb != nil {
				err := rdb.Publish(ctx, cfg.ErrorChannel, log).Err()
				if err != nil {
					s.logger.WithError(err).Error("Failed to publish log to Redis. Writing to file")
					helper.WriteToFile(fallbackFile, log) // Fallback to file logging
				}
			} else {
				// Redis unavailable, directly write to file
				helper.WriteToFile(fallbackFile, log)
			}
		}
	}()

	go func() {
		for log := range s.logDebugChan {
			s.logger.Debugln(string(log))

			if rdb != nil {
				err := rdb.Publish(ctx, cfg.DebugChannel, log).Err()
				if err != nil {
					s.logger.WithError(err).Error("Failed to publish log to Redis. Writing to file")
					helper.WriteToFile(fallbackFile, log) // Fallback to file logging
				}
			} else {
				// Redis unavailable, directly write to file
				helper.WriteToFile(fallbackFile, log)
			}
		}
	}()
}

func (s *service) Close() {
	// Close the channels
	close(s.logChan)
	close(s.logInfoChan)
	close(s.logWarnChan)
	close(s.logErrorChan)
	close(s.logDebugChan)
}
