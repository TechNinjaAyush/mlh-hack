package db

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToIncidents() (*mongo.Database, error) {

	err := godotenv.Load(".env")
	MONGO_URL := os.Getenv("MONGO_LOCAL_URl")

	if err != nil {
		return nil, fmt.Errorf("error loading .env: %v", err)
	}

	fmt.Printf("mongourl is %s", MONGO_URL)

	clientOptions := options.Client().ApplyURI(MONGO_URL)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("mongo connection failed: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("mongo ping failed: %v", err)
	}

	fmt.Println("Connection to db is successful...")

	db := client.Database("mlh")

	return db, nil
}
