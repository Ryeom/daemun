package cache

import (
	"context"
	logger "github.com/Ryeom/daemun/log"
	redisutil "github.com/Ryeom/daemun/redis"
	"github.com/go-redis/redis/v8"
)

// LimitPolicyCache : Redis에 저장된 limit policy들을 URL(및 조건 유형)별로 메모리에 캐싱
// EX. key "/endpoint1/발급:key-code" rule { "aaaa-1": 200 }
var LimitPolicyCache map[string]redisutil.LimitPolicyData

// LoadLimitPolicyCache : Redis에서 오늘 날짜(또는 전체 정책, 필요에 따라 필터 적용) 기준으로 정책들을 조회하여 캐싱
func LoadLimitPolicyCache(ctx context.Context, client *redis.Client) error {
	policies, err := redisutil.LoadTodayLimitPolicyConfigs(ctx, client)
	if err != nil {
		return err
	}

	policyCache := make(map[string]redisutil.LimitPolicyData)
	for key, policy := range policies {
		policyCache[key] = policy
	}
	LimitPolicyCache = policyCache
	logger.ServerLogger.Printf("캐시 로드 완료: %d개의 limit policy가 메모리에 저장됨", len(LimitPolicyCache))
	return nil
}
