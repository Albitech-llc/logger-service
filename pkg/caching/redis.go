package caching

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	//once sync.Once

	ctx context.Context
	rdb *redis.Client
	//err error
)

// InitializeRedis initializes the Redis client
// Return the client, context, and nil for error (indicating success)
func InitializeRedis(host string, port string, db int) (*redis.Client, context.Context, error) {
	//once.Do(func() {
	// Initialize Redis client
	redisAddr := fmt.Sprintf("%s:%s", host, port)
	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   db,
	})

	// Create a context
	ctx = context.Background()
	// Test the Redis connection
	err := rdb.Ping(ctx).Err()
	if err != nil {
		rdb = nil // Ensure rdb is nil if the connection fails
	}
	//})

	// Return the initialized Redis client, context, and any error encountered
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return rdb, ctx, nil
}
