package route

import (
	"github.com/Ryeom/daemun/cache"
	logger "github.com/Ryeom/daemun/log"
	redisutil "github.com/Ryeom/daemun/redis"
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"
)

// LimitPolicyHandler : 요청이 들어올 때 Redis에 캐시된 limit policy를 확인 -> 해당 정책의 제한값 이하에서만 요청을 처리
type LimitPolicyHandler struct {
	next        http.Handler
	redisClient *redis.Client
	manager     *redisutil.LimitManager
}

// NewLimitPolicyHandler : Redis 클라이언트를 주입받아 LimitPolicyHandler 인스턴스를 생성
func NewLimitPolicyHandler(redisClient *redis.Client, next http.Handler) http.Handler {
	return &LimitPolicyHandler{
		next:        next,
		redisClient: redisClient,
		manager:     redisutil.NewLimitManager(redisClient),
	}
}

// ServeHTTP : 요청 URL을 기준으로 cache에 저장된 정책을 조회 -> 정책의 제한값 이하에서만 요청을 처리
func (lph *LimitPolicyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path

	// 캐시에 저장된 정책을 조회
	policy, exists := cache.LimitPolicyCache[key]
	if !exists { // 없을 경우
		lph.next.ServeHTTP(w, r)
		return
	}

	limit := policy.DefaultLimit

	semKey := key
	ctx := r.Context()

	sem := lph.manager.GetSemaphore(semKey, limit)

	if err := sem.Acquire(ctx, 5*time.Second); err != nil {
		http.Error(w, "서버 과부하: 요청 제한 초과", http.StatusTooManyRequests)
		logger.ServerLogger.Printf("세마포어 획득 실패 (%s): %v", semKey, err)
		return
	}

	defer func() {
		if err := sem.Release(ctx); err != nil {
			logger.ServerLogger.Printf("세마포어 반환 오류 (%s): %v", semKey, err)
		}
		logger.TraceLogger.Printf("세마포어 반환 완료: %s", semKey)
	}()

	lph.next.ServeHTTP(w, r)
}
