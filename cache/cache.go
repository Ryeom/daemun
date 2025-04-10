package cache

import (
	"context"
	logger "github.com/Ryeom/daemun/log"
	redisutil "github.com/Ryeom/daemun/redis"
	"log"
)

func Init(ctx context.Context) error {
	err := redisutil.Init()
	if err != nil {
		logger.ServerLogger.Fatalf("Redis 연결 실패: %v", err)
	}

	client, err := redisutil.NewRedisClient(0)

	policies, err := redisutil.LoadTodayLimitPolicyConfigs(ctx, client)
	if err != nil {
		log.Printf("limit policy 캐시 로드 실패: %v", err)
		return err
	}

	LimitPolicyCache = policies
	logger.ServerLogger.Printf("limit policy 캐시 로드 완료: %d개의 정책이 메모리에 저장됨", len(policies))
	return nil
}
