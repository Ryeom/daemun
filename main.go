package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Ryeom/daemun/mongo"
	"github.com/Ryeom/daemun/route"
	"github.com/Ryeom/daemun/util/config"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.ConnectMongo(ctx)
	if err != nil {
		log.Fatalf("MongoDB 연결 실패: %v", err)
	}

	// config 설정
	appConfig, err := config.LoadConfig(ctx, client)
	if err != nil {
		log.Fatalf("설정 정보 로드 실패: %v", err)
	}

	// 로드된 설정을 기반으로 라우터 초기화
	router := route.NewRouter(appConfig)

	// HTTP 서버 시작
	addr := ":8080"
	log.Printf("서버 실행 중: %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("서버 에러: %v", err)
	}
}
