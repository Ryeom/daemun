package main

import (
	"context"
	logger "github.com/Ryeom/daemun/log"
	mongoutil "github.com/Ryeom/daemun/mongo"
	"github.com/Ryeom/daemun/route"
	"github.com/Ryeom/daemun/util/config"
	"github.com/gin-gonic/gin"
	"time"
)

func init() {
	if err := logger.Init(); err != nil {
		panic("로그 초기화 실패: " + err.Error())
	}
}
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	/* 1. mongo config load */
	err := mongoutil.Init()
	if err != nil {
		logger.ServerLogger.Fatalf("Mongo Connect 실패: %v", err)
	}
	/* 2. redis load */

	/* 3. config load */
	appConfig, err := config.LoadConfig(ctx)
	if err != nil {
		logger.ServerLogger.Fatalf("설정 정보 로드 실패: %v", err)
	}

	router := gin.Default()

	route.Initialize(router, appConfig)

	addr := ":8080"
	router.Run(addr)

}
