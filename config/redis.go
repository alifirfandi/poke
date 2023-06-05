package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var rdb *redis.Client

func InitRedis() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed connect to redis: %v", err)
	}

	return nil
}

func SetData(key, value string) error {
	err := rdb.Set(context.Background(), key, value, 10*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed insert data to redis: %v", err)
	}
	return nil
}

func GetData(key string) (string, error) {
	val, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		return "", fmt.Errorf("failed get data from redis: %v", err)
	}
	return val, nil
}
