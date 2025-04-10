package config

import (
	"context"
	mongoutil "github.com/Ryeom/daemun/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type RouteConfig struct { // 엔드포인트 설정 정보
	Key      string `bson:"key" json:"key"`           // 엔드포인트 이름
	IP       string `bson:"ip" json:"ip"`             // 대상 IP 또는 URL (필요에 따라 포트 포함)
	Platform string `bson:"platform" json:"platform"` // 플랫폼 따라 로드 정보 다름
}

type AppConfig struct { // 설정 정보
	Routes []RouteConfig `bson:"routes" json:"routes"`
}

func (ac *AppConfig) getRouteInfo(ctx context.Context) error {
	dbName := "daemun"
	collName := "endpoint_config"
	collection := mongoutil.Client.Database(dbName).Collection(collName)

	filter := bson.M{"platform": bson.M{"$in": []string{"common", "level5"}}}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(queryCtx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var routes []RouteConfig
	for cursor.Next(ctx) {
		var route RouteConfig
		if err := cursor.Decode(&route); err != nil {
			return err
		}
		routes = append(routes, route)
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	ac.Routes = routes
	return nil
}

func LoadConfig(ctx context.Context) (*AppConfig, error) {
	var err error
	appConfig := AppConfig{}

	/* 1. destination list */

	return &appConfig, err
}
