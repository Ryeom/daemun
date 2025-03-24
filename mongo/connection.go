package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo는 MongoDB에 연결하여 클라이언트를 반환합니다.
func ConnectMongo(ctx context.Context) (*mongo.Client, error) {
	// MongoDB 접속 URI (환경변수 등으로 관리할 수 있음)
	uri := "mongodb://localhost:27017"

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// 연결 확인을 위한 Ping
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
