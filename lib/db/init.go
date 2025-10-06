package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDB() *mongo.Database {
	// Get MongoDB credentials from environment variables
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")
	if mongoHost == "" {
		mongoHost = "localhost:27017"
	}
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, mongoPassword, mongoHost)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	return client.Database("health-tech")
}
