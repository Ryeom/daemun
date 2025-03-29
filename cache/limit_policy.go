package cache

import (
	"context"
	mg "github.com/Ryeom/daemun/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// LimitPolicyCache : 오늘 날짜 기준으로 활성화된 limit policy들을 URL을 키로 메모리에 캐싱
var LimitPolicyCache map[string]mg.LimitPolicyConfig

// LoadLimitPolicyCache : MongoDB에서 오늘 날짜에 해당하는 모든 정책들을 조회하여 캐싱
func LoadLimitPolicyCache(ctx context.Context, client *mongo.Client) error {
	policies, err := mg.LoadTodayLimitPolicyConfigs(ctx, client)
	if err != nil {
		return err
	}

	policyCache := make(map[string]mg.LimitPolicyConfig)
	for _, policy := range policies {
		policyCache[policy.URL] = policy
	}
	LimitPolicyCache = policyCache
	log.Printf("캐시 로드 완료: %d개의 limit policy가 메모리에 저장됨", len(LimitPolicyCache))
	return nil
}
