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
	var mongoUser, mongoPassword, mongoHost string
	mongoURI := os.Getenv("MONGO_URI")
	// Get MongoDB credentials from environment variables if in production environment
	if lib.IsProduction() {
		mongoHost = os.Getenv("MONGO_HOST")
		mongoUser = os.Getenv("MONGO_USER")
		mongoPassword = os.Getenv("MONGO_PASSWORD")
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
