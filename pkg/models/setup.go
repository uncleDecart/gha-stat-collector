package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		log.Println("MONGO_URL environment variable is not defined")
		return nil, fmt.Errorf("MONGO_URL environment variable is not defined")
	}
	mongoUser := os.Getenv("MONGO_USERNAME")
	if mongoUser == "" {
		log.Println("MONGO_USERNAME environment variable is not defined")
		return nil, fmt.Errorf("MONGO_USERNAME environment variable is not defined")
	}
	mongoPwd := os.Getenv("MONGO_PASSWORD")
	if mongoPwd == "" {
		log.Println("MONGO_PASSWORD environment variable is not defined")
		return nil, fmt.Errorf("MONGO_PASSWORD environment variable is not defined")
	}

	credential := options.Credential{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
	}
	clientOptions := options.Client().ApplyURI(mongoURL).SetAuth(credential)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Error connecting to database ", err)
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("Error pining database ", err)
		return nil, err
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return client, nil
}
