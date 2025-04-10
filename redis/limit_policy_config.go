package redisutil

import (
	"context"
	"encoding/json"

	logger "github.com/Ryeom/daemun/log"
	"github.com/go-redis/redis/v8"
)

type LimitPolicyData struct {
	DefaultLimit int            `json:"default_limit"` // 기본 제한값
	Rules        map[string]int `json:"rules"`         // 조건 값을 key로, 제한값을 value로 (예: {"aaaa-1":200})
}

const RedisLimitPolicyKey = "limit_policy_config"

// SaveLimitPolicies는 주어진 policies 맵을 Redis 해시에 저장합니다.
func SaveLimitPolicies(ctx context.Context, client *redis.Client, policies map[string]LimitPolicyData) error {
	data := make(map[string]interface{})
	for key, policy := range policies {
		bytes, err := json.Marshal(policy)
		if err != nil {
			logger.ServerLogger.Printf("JSON marshal 실패 for %s: %v", key, err)
			continue
		}
		data[key] = string(bytes)
	}
	if err := client.HSet(ctx, RedisLimitPolicyKey, data).Err(); err != nil {
		logger.ServerLogger.Printf("Redis HSet 실패: %v", err)
		return err
	}
	logger.ServerLogger.Printf("limit_policy_config 저장 완료, 정책 수: %d", len(data))
	return nil
}

// LoadLimitPolicies는 Redis 해시에서 모든 정책 데이터를 조회하여 반환합니다.
func LoadLimitPolicies(ctx context.Context, client *redis.Client) (map[string]LimitPolicyData, error) {
	result, err := client.HGetAll(ctx, RedisLimitPolicyKey).Result()
	if err != nil {
		logger.ServerLogger.Printf("Redis HGetAll 실패: %v", err)
		return nil, err
	}
	policies := make(map[string]LimitPolicyData)
	for key, jsonStr := range result {
		var policy LimitPolicyData
		if err := json.Unmarshal([]byte(jsonStr), &policy); err != nil {
			logger.ServerLogger.Printf("JSON unmarshal 실패 for %s: %v", key, err)
			continue
		}
		policies[key] = policy
	}
	logger.ServerLogger.Printf("limit_policy_config 로드 완료, 정책 수: %d", len(policies))
	return policies, nil
}

// LoadTodayLimitPolicyConfigs는 오늘 날짜(혹은 전체 정책, 필요 시 필터 적용) 기준의 정책들을 반환합니다.
// 여기서는 간단하게 모든 정책을 반환합니다.
func LoadTodayLimitPolicyConfigs(ctx context.Context, client *redis.Client) (map[string]LimitPolicyData, error) {
	return LoadLimitPolicies(ctx, client)
}
