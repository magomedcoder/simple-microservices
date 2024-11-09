package internal

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var client *mongo.Client

type Config struct {
	Models Models
}

func ConnectToMongo() (*mongo.Client, error) {
	dsn := os.Getenv("DSN_MONGODB")
	clientOptions := options.Client().ApplyURI(dsn)
	clientOptions.SetAuth(options.Credential{
		Username: "mongo",
		Password: "mongo",
	})
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Ошибка подключения:", err)
		return nil, err
	}
	log.Println("Подключено к MongoDB")
	return c, nil
}
