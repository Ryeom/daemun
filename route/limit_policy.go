package route

import (
	"encoding/json"
	"net/http"
	"time"

	logger "github.com/Ryeom/daemun/log"
	"github.com/gin-gonic/gin"

	redisutil "github.com/Ryeom/daemun/redis"
	"github.com/go-redis/redis/v8"
)

// GinDistributedLimitMiddleware는 Redis 해시 "limit_policy_config"에서 해당 엔드포인트 정책을 직접 조회하여,
// 분산 세마포어를 통해 요청 제한을 적용하는 Gin 미들웨어입니다.
func GinDistributedLimitMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	// 분산 제한 매니저 생성
	manager := redisutil.NewLimitManager(redisClient)

	return func(c *gin.Context) {
		// 엔드포인트 키는 요청 URL 경로를 사용 (예: "/endpoint1/발급:key-code")
		semKey := c.Request.URL.Path
		ctx := c.Request.Context()
		redisKey := redisutil.RedisLimitPolicyKey

		// Redis에서 해당 키에 대한 정책 JSON 문자열 조회
		policyJSON, err := redisClient.HGet(ctx, redisKey, semKey).Result()
		if err == redis.Nil {
			// 정책이 없으면 제한 없이 처리
			c.Next()
			return
		} else if err != nil {
			logger.ServerLogger.Printf("Redis HGet 실패 (%s): %v", semKey, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 조회 오류"})
			c.Abort()
			return
		}

		var policy redisutil.LimitPolicyData
		if err := json.Unmarshal([]byte(policyJSON), &policy); err != nil {
			logger.ServerLogger.Printf("JSON 파싱 오류 (%s): %v", semKey, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 파싱 오류"})
			c.Abort()
			return
		}

		// 기본 제한값 적용
		limit := policy.DefaultLimit

		// Redis 기반 세마포어 획득
		sem := manager.GetSemaphore(semKey, limit)
		if err := sem.Acquire(ctx, 5*time.Second); err != nil {
			c.String(http.StatusTooManyRequests, "서버 과부하: 요청 제한 초과")
			logger.ServerLogger.Printf("세마포어 획득 실패 (%s): %v", semKey, err)
			c.Abort()
			return
		}

		defer func() {
			if err := sem.Release(ctx); err != nil {
				logger.ServerLogger.Printf("세마포어 반환 오류 (%s): %v", semKey, err)
			}
			logger.TraceLogger.Printf("세마포어 반환 완료: %s", semKey)
		}()

		c.Next()
	}
}
