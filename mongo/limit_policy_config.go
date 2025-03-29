package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

// LimitPolicyConfig : 하나의 URL 및 접근 제한 정책 문서
type LimitPolicyConfig struct {
	URL          string `bson:"url" json:"url"`                     // 예: "ticket/purchase/reservation" 혹은 전체 URL
	DefaultLimit int    `bson:"default_limit" json:"default_limit"` // 기본 제한 값
	Period       Period `bson:"period" json:"period"`               // 정책 적용 기간
	OrderID      string `bson:"order_id" json:"order_id"`           // 정책 순서나 식별용 ID
	Rules        []Rule `bson:"rules" json:"rules"`                 // 조건별 제한 정책 목록
}

// Period : 정책이 유효한 기간
type Period struct {
	Start time.Time `bson:"start" json:"start"`
	End   time.Time `bson:"end" json:"end"`
}

// Rule : 특정 조건에 따른 제한 정책
// 예시 : { "match": { "header": {"X-Event-Name": "캣츠 무비"} }, "limit": 500 }
type Rule struct {
	Match map[string]map[string]string `bson:"match" json:"match"`
	Limit int                          `bson:"limit" json:"limit"`
}

// getLimitPolicyConfigCollection : limit_policy_config 컬렉션을 반환
func getLimitPolicyConfigCollection(client *mongo.Client) *mongo.Collection {
	return client.Database("daemun").Collection("limit_policy_config")
}

// CreateLimitPolicyConfig : 새로운 접근 제한 정책 문서를 생성
func CreateLimitPolicyConfig(ctx context.Context, client *mongo.Client, policy LimitPolicyConfig) error {
	collection := getLimitPolicyConfigCollection(client)
	_, err := collection.InsertOne(ctx, policy)
	if err != nil {
		log.Printf("문서 삽입 실패 (url: %s): %v", policy.URL, err)
	}
	return err
}

// GetLimitPolicyConfigByURL : URL을 기준으로 단일 접근 제한 정책 문서를 조회
func GetLimitPolicyConfigByURL(ctx context.Context, client *mongo.Client, url string) (*LimitPolicyConfig, error) {
	collection := getLimitPolicyConfigCollection(client)
	filter := bson.M{"url": url}
	var policy LimitPolicyConfig
	err := collection.FindOne(ctx, filter).Decode(&policy)
	if err != nil {
		log.Printf("문서 조회 실패 (url: %s): %v", url, err)
		return nil, err
	}
	return &policy, nil
}

// UpdateLimitPolicyConfig : URL을 기준으로 접근 제한 정책 문서를 수정
// update 파라미터에는 수정할 필드와 값을 담은 bson.M을 전달
func UpdateLimitPolicyConfig(ctx context.Context, client *mongo.Client, url string, update bson.M) error {
	collection := getLimitPolicyConfigCollection(client)
	filter := bson.M{"url": url}
	updateDoc := bson.M{"$set": update}
	result, err := collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		log.Printf("문서 수정 실패 (url: %s): %v", url, err)
		return err
	}
	log.Printf("문서 수정 성공 (url: %s): %d건 수정됨", url, result.ModifiedCount)
	return nil
}

// DeleteLimitPolicyConfig : URL을 기준으로 접근 제한 정책 문서를 삭제
func DeleteLimitPolicyConfig(ctx context.Context, client *mongo.Client, url string) error {
	collection := getLimitPolicyConfigCollection(client)
	filter := bson.M{"url": url}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("문서 삭제 실패 (url: %s): %v", url, err)
		return err
	}
	log.Printf("문서 삭제 성공 (url: %s): %d건 삭제됨", url, result.DeletedCount)
	return nil
}
