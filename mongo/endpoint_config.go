package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

// EndpointConfig : API Gateway에서 사용되는 엔드포인트 설정 정보
type EndpointConfig struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Key       string             `bson:"key" json:"key"`
	IP        string             `bson:"ip" json:"ip"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func getEndpointConfigCollection(client *mongo.Client) *mongo.Collection {
	return client.Database("daemun").Collection("endpoint_config")
}

func GetAllEndpointConfigs(ctx context.Context, client *mongo.Client) ([]EndpointConfig, error) {
	collection := getEndpointConfigCollection(client)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var configs []EndpointConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func CreateEndpointConfig(ctx context.Context, client *mongo.Client, config EndpointConfig) error {
	collection := getEndpointConfigCollection(client)
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()
	_, err := collection.InsertOne(ctx, config)
	if err != nil {
		log.Printf("EndpointConfig 생성 실패 (key: %s): %v", config.Key, err)
	}
	return err
}

func UpdateEndpointConfig(ctx context.Context, client *mongo.Client, id string, updateData bson.M) error {
	collection := getEndpointConfigCollection(client)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	updateData["updated_at"] = time.Now()
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	return err
}

func DeleteEndpointConfig(ctx context.Context, client *mongo.Client, id string) error {
	collection := getEndpointConfigCollection(client)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
