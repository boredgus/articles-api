package db

import (
	"context"
	"fmt"
	"sync"
	"time"
	"user-management/config"
	"user-management/internal/gateways"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisClient struct {
	*redis.Client
	ExpiresAfter time.Duration
}

var redisClient *RedisClient
var redisOnce sync.Once

func NewRedisStore() gateways.CacheStore {
	redisOnce.Do(func() {
		cfg := config.GetConfig()
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%v:%v", cfg.RedisContainer, cfg.RedisPort),
			Password: cfg.RedisPassword,
		})
		_, err := client.Ping(context.Background()).Result()
		if err != nil {
			logrus.Fatal("failed to connect to redis", err)
		}
		redisClient = &RedisClient{
			Client:       client,
			ExpiresAfter: time.Duration(cfg.RedisExpirationMinutes) * time.Minute}
	})
	return &RedisStore{client: redisClient}
}

type RedisStore struct {
	client *RedisClient
}

func (s *RedisStore) Get(key string, value any) error {
	logrus.Infoln("> redis.Get:")
	return s.client.Get(context.Background(), key).Scan(value)
}

func (s *RedisStore) Set(key string, value any) error {
	logrus.Infoln("> redis.Set:")
	return s.client.Set(context.Background(), key, value, s.client.ExpiresAfter).Err()
}
