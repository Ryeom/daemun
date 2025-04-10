package redisutil

import (
	"sync"

	"github.com/go-redis/redis/v8"
)

type LimitManager struct {
	redisClient *redis.Client
	semaphores  map[string]*DistributedSemaphore
	mu          sync.Mutex
}

func NewLimitManager(redisClient *redis.Client) *LimitManager {
	return &LimitManager{
		redisClient: redisClient,
		semaphores:  make(map[string]*DistributedSemaphore),
	}
}

// GetSemaphore는 주어진 key에 대해 기존의 DistributedSemaphore가 있으면 반환하고,
// 없으면 maxTokens(limit)로 새로 생성하여 저장한 후 반환합니다.
func (lm *LimitManager) GetSemaphore(key string, limit int) *DistributedSemaphore {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	if sem, ok := lm.semaphores[key]; ok {
		return sem
	}
	sem := NewDistributedSemaphore(lm.redisClient, key, limit)
	lm.semaphores[key] = sem
	return sem
}
