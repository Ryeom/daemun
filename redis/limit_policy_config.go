package redisutil

import (
	"context"
	"encoding/json"
	logger "github.com/Ryeom/daemun/log"
	"github.com/go-redis/redis/v8"
)

// LimitPolicyData : 각 URL(및 조건 유형)에 대해 저장할 정책 데이터
// "/endpoint1/발급:key-code": { "aaaa-1": 200 }
type LimitPolicyData struct {
	DefaultLimit int            `json:"default_limit"`
	Rules        map[string]int `json:"rules"`
}

// SaveLimitPolicies : policies 맵(각 엔드포인트 키에 해당하는 LimitPolicyData)을 Redis의 해시("limit_policy_config")
func SaveLimitPolicies(ctx context.Context, client *redis.Client, policies map[string]LimitPolicyData) error {
	redisKey := "limit_policy_config"
	data := make(map[string]interface{})
	for key, policy := range policies {
		bytes, err := json.Marshal(policy)
		if err != nil {
			logger.ServerLogger.Printf("JSON marshal 실패 for %s: %v", key, err)
			continue
		}
		data[key] = string(bytes)
	}
	if err := client.HSet(ctx, redisKey, data).Err(); err != nil {
		logger.ServerLogger.Printf("Redis HSet 실패: %v", err)
		return err
	}
	logger.ServerLogger.Printf("limit_policy_config 저장 완료, 정책 수: %d", len(data))
	return nil
}

// LoadLimitPolicies : Redis의 해시("limit_policy_config")에서 모든 정책 데이터를 읽어와 맵으로 반환
func LoadLimitPolicies(ctx context.Context, client *redis.Client) (map[string]LimitPolicyData, error) {
	redisKey := "limit_policy_config"
	result, err := client.HGetAll(ctx, redisKey).Result()
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

// LoadTodayLimitPolicyConfigs : 오늘 날짜(오늘 00:00 ~ 내일 00:00)에 유효한 정책들을 조회하여 반환
func LoadTodayLimitPolicyConfigs(ctx context.Context, client *redis.Client) (map[string]LimitPolicyData, error) {
	return LoadLimitPolicies(ctx, client)
}
