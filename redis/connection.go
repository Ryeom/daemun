package redisutil

import (
	"context"
	"time"

	logger "github.com/Ryeom/daemun/log"
	"github.com/go-redis/redis/v8"
)

var LimitPolicy *redis.Client

func Init() error {
	var err error
	LimitPolicy, err = NewRedisClient(0)
	if err != nil {
		logger.ServerLogger.Fatalln("limitPolicy Redis 연결 실패: %v", err)
	}
	return err
}

func NewRedisClient(num int) (*redis.Client, error) {
	options := &redis.Options{
		Addr:         "",
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
	err := client.Ping(ctx).Err()
	if err != nil {
		logger.ServerLogger.Printf("Redis 연결 실패: %v", err)
		// TODO 재시도 로직 or 패닉 처리
	} else {
		logger.ServerLogger.Printf("Redis에 성공적으로 연결됨: %s", options.Addr)
	}
	return client, err
}
