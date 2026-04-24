package db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client
var Ctx = context.Background()

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // 没有密码填空
		DB:       0,                // 默认使用 0 号库
	})

	// 测试连接
	_, err := Redis.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Redis 连接失败 (请确认 Redis 服务已启动): %v", err)
	}
	log.Println("✅ Redis 初始化成功")
}
