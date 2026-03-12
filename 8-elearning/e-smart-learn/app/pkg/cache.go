package pkg

import (
	"context"
	"fmt"
	"time"

	"elearning-api/config"

	"github.com/redis/go-redis/v9"
)

type CacheProvider interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	Delete(ctx context.Context, key string) error
}

type cacheProvider struct {
	client *redis.Client
}

func NewCacheProvider(redisConfig *config.RedisConfig) CacheProvider {
	fmt.Printf("Connecting to Redis at %s:%d...\n", redisConfig.Host, redisConfig.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil
	}
	return &cacheProvider{
		client: rdb,
	}
}

func (c *cacheProvider) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *cacheProvider) Get(ctx context.Context, key string) (any, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *cacheProvider) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
