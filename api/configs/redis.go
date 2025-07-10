package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

func ConnectRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDBString := os.Getenv("REDIS_DB")

	if redisAddr == "" {
		log.Fatal("REDIS_ADDR is not set in .env")
	}

	redisDB, err := strconv.Atoi(redisDBString)
	if err != nil {
		log.Printf("Warning: Redis DB is not a valid number, , defaulting to DB 0: %v", err)
		redisDB = 0
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	res, err := RedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatalf("Could not connect to redis: %v", err)

	}
	fmt.Println("Connected to Redis:", res)
}

func CloseRedis() {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			log.Printf("Error closing redis client: %v", err)
		} else {
			fmt.Println("Redis client closed")
		}
	}
}

func GetRedisClient() *redis.Client {
	return RedisClient
}
