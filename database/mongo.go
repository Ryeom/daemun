package database

import (
	"context"
	"fmt"
	"github.com/Ryeom/daemun/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func New(platform, target, port string) *mongo.Client {
	ip := ""
	if false {
		log.Logger.Error(ip + ":" + port + " 통신 불가.")
		return nil
	}
	return newMongoClient(ip)
}

func newMongoClient(key string) *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	clientOptions := options.Client().ApplyURI(key).SetMaxPoolSize(3)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Logger.Error("Client Connection %s", err)

	}
	client.Ping(ctx, nil)
	if err != nil {
		log.Logger.Error("Client Ping %s", err)

	}
	return client
}

func SelectAll(client *mongo.Client, where map[string]string) map[string]string {
	result := map[string]string{}
	var l []bson.E
	for i, v := range where {
		l = append(l, bson.E{Key: i, Value: v})
	}
	//E의 배열이 D
	collection := client.Database("gateway_configuration").Collection("endpoint")
	cursor, err := collection.Find(context.TODO(), bson.D(l))
	if err != nil {
		log.Logger.Error("Find %s", err)
	}

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Logger.Error("All %s", err)
	}

	for _, v := range results {
		fmt.Println(v)
	}
	return result
}
