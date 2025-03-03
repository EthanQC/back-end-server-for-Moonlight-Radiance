// Redis缓存连接

package common

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis 初始化Redis客户端
func InitRedis(addr string, password string, db int) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password by default
		DB:       db,       // use default DB
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("could not connect to redis: %v", err)
	}
	log.Println("Redis connected successfully.")
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	err := RedisClient.Close()
	if err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}
}
