package caching

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	once sync.Once

	ctx context.Context
	RDB *redis.Client
)

// InitializeRedis initializes the Redis client
// Return the client, context, and nil for error (indicating success)
func InitializeRedis(host string, port string, db int) (*redis.Client, context.Context, error) {
	once.Do(func() {
		// Initialize Redis client
		redisAddr := fmt.Sprintf("%s:%s", host, port)
		RDB = redis.NewClient(&redis.Options{
			Addr: redisAddr,
			DB:   db,
		})

		// Create a context
		ctx = context.Background()
		// Test the Redis connection
		err := RDB.Ping(ctx).Err()
		if err != nil {
			log.Printf("failed to connect to Redis: %v\n", err)
		} else {
			fmt.Println("Redis connection initialized successfully")
		}
	})

	return RDB, ctx, nil
}
