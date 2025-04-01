package redisutil

import (
	"context"
	"fmt"
	"time"

	logger "github.com/Ryeom/daemun/log"
	"github.com/go-redis/redis/v8"
)

// DistributedSemaphore : Redis를 이용한 분산 세마포어를 구현. 각 세마포어는 하나의 고유 키와 최대 토큰 수(maxTokens)를 관리
type DistributedSemaphore struct {
	redisClient *redis.Client
	key         string
	maxTokens   int
}

// luaAcquire : Redis Lua 스크립트를 이용해 토큰을 원자적으로 획득
// 스크립트는 지정한 키의 현재 값이 0보다 크면 1을 감소시키고 1을 반환
var luaAcquire = redis.NewScript(`
	local current = tonumber(redis.call("GET", KEYS[1]) or "0")
	if current > 0 then
		redis.call("DECR", KEYS[1])
		return 1
	else
		return 0
	end
`)

// NewDistributedSemaphore : 주어진 Redis 클라이언트와 키, 최대 토큰 수를 이용해 DistributedSemaphore 생성 : 초기화 시, 해당 키가 존재하지 않으면 maxTokens 값으로 초기화
func NewDistributedSemaphore(redisClient *redis.Client, key string, maxTokens int) *DistributedSemaphore {
	ctx := context.Background()
	err := redisClient.SetNX(ctx, key, maxTokens, 0).Err()
	if err != nil {
		logger.ServerLogger.Printf("DistributedSemaphore 초기화 오류 (%s): %v", key, err)
	}
	return &DistributedSemaphore{
		redisClient: redisClient,
		key:         key,
		maxTokens:   maxTokens,
	}
}

// Acquire : Lua 스크립트를 이용해 토큰 1 획득
func (ds *DistributedSemaphore) Acquire(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		res, err := luaAcquire.Run(ctx, ds.redisClient, []string{ds.key}).Result()
		if err != nil {
			return err
		}
		if res.(int64) == 1 {
			logger.TraceLogger.Printf("DistributedSemaphore acquired: %s", ds.key)
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout acquiring semaphore %s", ds.key)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// Release : Lua 스크립트를 이용해 토큰을 반환
func (ds *DistributedSemaphore) Release(ctx context.Context) error {
	luaRelease := redis.NewScript(`
		local current = tonumber(redis.call("GET", KEYS[1]) or "0")
		if current < tonumber(ARGV[1]) then
			return redis.call("INCR", KEYS[1])
		else
			return current
		end
	`)
	res, err := luaRelease.Run(ctx, ds.redisClient, []string{ds.key}, ds.maxTokens).Result()
	if err != nil {
		return err
	}
	logger.TraceLogger.Printf("DistributedSemaphore released: %s, new value: %v", ds.key, res)
	return nil
}
