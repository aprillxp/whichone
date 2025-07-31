package services

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		panic("Couldn't connect to redis: " + err.Error())
	}
}

func AcquireLock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	success, err := RDB.SetNX(ctx, key, "locked", expiration).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

func ReleaseLock(ctx context.Context, key string) error {
	_, err := RDB.Del(ctx, key).Result()
	return err
}
