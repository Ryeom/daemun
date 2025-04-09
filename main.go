package main

import (
	"context"
	"github.com/Ryeom/daemun/cache"
	logger "github.com/Ryeom/daemun/log"
	"github.com/Ryeom/daemun/mongo"
	"github.com/Ryeom/daemun/route"

	"github.com/Ryeom/daemun/util/config"
	"net/http"
	"time"
)

func main() {
	if err := logger.Init(); err != nil {
		panic("로그 초기화 실패: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.ConnectMongo(ctx)
	if err != nil {
		logger.ServerLogger.Fatalf("MongoDB 연결 실패: %v", err)
	}

	appConfig, err := config.LoadConfig(ctx, client)
	if err != nil {
		logger.ServerLogger.Fatalf("설정 정보 로드 실패: %v", err)
	}

	if err = cache.Init(ctx); err != nil {
		logger.ServerLogger.Fatalf("캐시 초기화 실패: %v", err)
	}

	router := route.NewRouter(appConfig)

	addr := ":8080"
	logger.ServerLogger.Printf("서버 실행 중: %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.ServerLogger.Fatalf("서버 에러: %v", err)
	}
}
