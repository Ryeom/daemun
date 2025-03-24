package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RouteConfig struct { // 엔드포인트 설정 정보
	Key      string `bson:"key" json:"key"`           // 엔드포인트 이름
	IP       string `bson:"ip" json:"ip"`             // 대상 IP 또는 URL (필요에 따라 포트 포함)
	Platform string `bson:"platform" json:"platform"` // 플랫폼 따라 로드 정보 다름
}

type AppConfig struct { // 설정 정보
	Routes []RouteConfig `bson:"routes" json:"routes"`
}

func LoadConfig(ctx context.Context, client *mongo.Client) (*AppConfig, error) {
	dbName := "daemun"
	collName := "endpoint_config"
	collection := client.Database(dbName).Collection(collName)

	filter := bson.M{"platform": bson.M{"$in": []string{"common", "level5"}}}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(queryCtx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var routes []RouteConfig
	for cursor.Next(ctx) {
		var route RouteConfig
		if err := cursor.Decode(&route); err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	config := AppConfig{Routes: routes}
	log.Printf("로드된 설정: %+v", config)
	return &config, nil
}
