package redisutil

import (
	"context"
	"time"

	logger "github.com/Ryeom/daemun/log"
	"github.com/go-redis/redis/v8"
)

var Client map[int]*redis.Client

const (
	limitPolicy = 1
)

func init() {
	Client[limitPolicy] = NewRedisClient(0)
}

func NewRedisClient(num int) *redis.Client {
	options := &redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		DB:           num,
		PoolSize:     10,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		logger.ServerLogger.Printf("Redis 연결 실패: %v", err)
		// TODO 재시도 로직 or 패닉 처리
	} else {
		logger.ServerLogger.Printf("Redis에 성공적으로 연결됨: %s", options.Addr)
	}
	return client
}
