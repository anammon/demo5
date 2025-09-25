package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

func InitRedis() {
	// 初始化 Redis 连接
	redisclient := redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Addr,
		Password: AppConfig.Redis.PassWord, // no password set
		DB:       AppConfig.Redis.DB,       // use default DB
	})
	if _, err := redisclient.Ping(RedisCtx).Result(); err != nil {
		log.Fatalf("unable to connect redis-%v", err) //记录日志
	}
	RedisClient = redisclient
}
