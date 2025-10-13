package db

import (
	"context"
	"fmt"
	"health-tech/lib"
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
		mongoHost = "localhost:27017/?replicaSet=rs0"
	}
	mongoURI := "mongodb://localhost:27017/?replicaSet=rs0"
	if lib.IsProduction() {
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, mongoPassword, mongoHost)
	}
	clientOptions := options.Client().ApplyURI(mongoURI)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	return client.Database("health-tech")
}
