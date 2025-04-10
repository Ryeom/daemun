package route

import (
	"context"
	"net/http"
	"time"

	redisutil "github.com/Ryeom/daemun/redis"
	"github.com/gin-gonic/gin"
)

type LimitPolicyRequest struct {
	Key          string         `json:"key"`
	DefaultLimit int            `json:"default_limit"`
	Rules        map[string]int `json:"rules"`
}

func GetLimitPolicyHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	redisClient := redisutil.LimitPolicy
	policies, err := redisutil.LoadLimitPolicies(ctx, redisClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 로드 실패: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, policies)
}

func CreateLimitPolicyHandler(c *gin.Context) {
	var req LimitPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	policy := redisutil.LimitPolicyData{
		DefaultLimit: req.DefaultLimit,
		Rules:        req.Rules,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	redisClient := redisutil.LimitPolicy

	policies, err := redisutil.LoadLimitPolicies(ctx, redisClient)
	if err != nil {
		policies = make(map[string]redisutil.LimitPolicyData)
	}

	policies[req.Key] = policy

	if err := redisutil.SaveLimitPolicies(ctx, redisClient, policies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 저장 실패: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "정책 생성됨", "policy": req})
}

func UpdateLimitPolicyHandler(c *gin.Context) {
	var req LimitPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	redisClient := redisutil.LimitPolicy

	policies, err := redisutil.LoadLimitPolicies(ctx, redisClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 로드 실패: " + err.Error()})
		return
	}

	if _, exists := policies[req.Key]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "정책이 존재하지 않습니다."})
		return
	}

	policies[req.Key] = redisutil.LimitPolicyData{
		DefaultLimit: req.DefaultLimit,
		Rules:        req.Rules,
	}

	if err := redisutil.SaveLimitPolicies(ctx, redisClient, policies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 업데이트 실패: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "정책 업데이트됨", "policy": req})
}

func DeleteLimitPolicyHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "삭제할 정책의 key가 필요합니다."})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	redisClient := redisutil.LimitPolicy

	policies, err := redisutil.LoadLimitPolicies(ctx, redisClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 로드 실패: " + err.Error()})
		return
	}

	if _, exists := policies[key]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "정책이 존재하지 않습니다."})
		return
	}

	delete(policies, key)

	if err := redisutil.SaveLimitPolicies(ctx, redisClient, policies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 삭제 실패: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "정책 삭제됨", "key": key})
}
