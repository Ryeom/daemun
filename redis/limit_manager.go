package redisutil

import (
	"sync"

	"github.com/go-redis/redis/v8"
)

// LimitManager : 분산 세마포어들을 관리하는 매니저 각 고유 키에 대해 DistributedSemaphore를 생성 및 재사용
type LimitManager struct {
	redisClient *redis.Client
	semaphores  map[string]*DistributedSemaphore
	mu          sync.Mutex
}

// NewLimitManager : 주어진 Redis 클라이언트를 사용해 새로운 LimitManager를 생성
func NewLimitManager(redisClient *redis.Client) *LimitManager {
	return &LimitManager{
		redisClient: redisClient,
		semaphores:  make(map[string]*DistributedSemaphore),
	}
}

// GetSemaphore : 주어진 key에 대해 기존 세마포어가 있으면 반환 / 없으면 주어진 제한(limit)으로 새로 생성한 후 저장하여 반환
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
