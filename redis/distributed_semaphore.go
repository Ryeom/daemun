package redisutil

import (
	"context"
	"fmt"
	"time"

	logger "github.com/Ryeom/daemun/log"
	"github.com/go-redis/redis/v8"
)

// Lua 스크립트를 이용해 토큰 획득(원자적 DECR)
var luaAcquire = redis.NewScript(`
	local current = tonumber(redis.call("GET", KEYS[1]) or "0")
	if current > 0 then
		redis.call("DECR", KEYS[1])
		return 1
	else
		return 0
	end
`)

// Lua 스크립트를 이용해 토큰 반환(원자적 INCR)
var luaRelease = redis.NewScript(`
	local current = tonumber(redis.call("GET", KEYS[1]) or "0")
	if current < tonumber(ARGV[1]) then
		return redis.call("INCR", KEYS[1])
	else
		return current
	end
`)

// DistributedSemaphore는 Redis를 이용해 분산 세마포어를 구현한 구조체입니다.
type DistributedSemaphore struct {
	redisClient *redis.Client
	key         string // 예: "semaphore:/endpoint:aaaa-1"
	maxTokens   int
}

// NewDistributedSemaphore는 지정한 Redis 클라이언트와 키, 최대 토큰 수를 사용해 세마포어를 생성합니다.
// 초기화 시, 해당 키가 없으면 maxTokens로 설정합니다.
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

// Acquire는 지정한 timeout 내에 Lua 스크립트를 통해 토큰을 획득합니다.
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

// Release는 Lua 스크립트를 통해 토큰을 반환합니다.
func (ds *DistributedSemaphore) Release(ctx context.Context) error {
	res, err := luaRelease.Run(ctx, ds.redisClient, []string{ds.key}, ds.maxTokens).Result()
	if err != nil {
		return err
	}
	logger.TraceLogger.Printf("DistributedSemaphore released: %s, new value: %v", ds.key, res)
	return nil
}
