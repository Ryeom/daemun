package cache

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"

	"log"
)

// Init : 현재 필요한 모든 캐시 데이터를 초기화
func Init(ctx context.Context, client *mongo.Client) error {
	if err := LoadLimitPolicyCache(ctx, client); err != nil {
		log.Printf("limit policy 캐시 로드 실패: %v", err)
		return err
	}

	// TODO: 다른 캐시 데이터 초기화 (예: 사용자 세션, 설정값 등)

	log.Printf("모든 캐시 초기화 완료")
	return nil
}
