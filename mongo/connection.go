package mongo

import (
	"context"
	"fmt"
	logger "github.com/Ryeom/daemun/log"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := getMongoUri()
	c, err := ConnectMongo(ctx, uri)
	if err != nil {
		logger.ServerLogger.Println("Mongo Connect 실패:", err)
		return err
	}
	Client = c
	return nil
}

func getMongoUri() string {
	username := ""
	password := url.QueryEscape("")
	host := ""
	port := ""
	dbName := ""

	authSource := "users"
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s",
		username, password, host, port, dbName, authSource)
	fmt.Println(uri)
	return "mongodb://localhost:27017"
}

func ConnectMongo(ctx context.Context, uri string) (*mongo.Client, error) {
	fmt.Println("MongoDB URI:", uri)
	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
