package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	logger "github.com/Ryeom/daemun/log"
	redisutil "github.com/Ryeom/daemun/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type LimitPolicyData = redisutil.LimitPolicyData

func GinDistributedLimitMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	manager := redisutil.NewLimitManager(redisClient)

	return func(c *gin.Context) {
		semKey := c.Request.URL.Path

		ctx := c.Request.Context()
		redisKey := "limit_policy_config"
		jsonStr, err := redisClient.HGet(ctx, redisKey, semKey).Result()
		if errors.Is(err, redis.Nil) {
			c.Next()
			return
		} else if err != nil {
			logger.ServerLogger.Printf("Redis HGet 실패 (%s): %v", semKey, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 조회 오류"})
			c.Abort()
			return
		}

		var policy LimitPolicyData
		if err := json.Unmarshal([]byte(jsonStr), &policy); err != nil {
			logger.ServerLogger.Printf("JSON 파싱 오류 (%s): %v", semKey, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 파싱 오류"})
			c.Abort()
			return
		}

		limit := policy.DefaultLimit

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
