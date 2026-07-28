package redis

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"time"

	"github.com/dmi3midd/grpcsso/internal/config"

	"github.com/redis/go-redis/v9"
)

// RedisService represents a service that interacts with Redis.
type RedisService interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the Redis client connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// GetClient returns the Redis client instance.
	GetClient() *redis.Client
}

type redisService struct {
	cfg *config.RedisConfig
	rdb *redis.Client
}

// New creates and initializes a new RedisService.
func New(cfg *config.RedisConfig) (RedisService, error) {
	opt, err := redis.ParseURL(cfg.URI)
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(opt)

	// Ping the redis server to ensure the connection is active
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	return &redisService{
		cfg: cfg,
		rdb: rdb,
	}, nil
}

// Health checks the health of the Redis connection by pinging it.
// It returns a map with keys indicating various health statistics.
func (s *redisService) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping Redis
	err := s.rdb.Ping(ctx).Err()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("redis down: %v", err)
		slog.Warn("redis down", slog.String("error", err.Error()))
		return stats
	}

	// Redis is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get pool statistics
	poolStats := s.rdb.PoolStats()
	stats["total_connections"] = strconv.Itoa(int(poolStats.TotalConns))
	stats["idle_connections"] = strconv.Itoa(int(poolStats.IdleConns))
	stats["hits"] = strconv.Itoa(int(poolStats.Hits))
	stats["misses"] = strconv.Itoa(int(poolStats.Misses))
	stats["timeouts"] = strconv.Itoa(int(poolStats.Timeouts))

	return stats
}

// Close closes the Redis client connection.
// It logs a message indicating the disconnection from Redis.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *redisService) Close() error {
	slog.Info("Disconnected from redis", slog.String("uri", s.cfg.URI))
	return s.rdb.Close()
}

// GetClient returns the Redis client instance.
func (s *redisService) GetClient() *redis.Client {
	return s.rdb
}
