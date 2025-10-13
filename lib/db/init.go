package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDB() *mongo.Database {
	// Get MongoDB credentials from environment variables
	/* mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")
	if mongoHost == "" {
		mongoHost = "localhost:27017/health-tech?authSource=admin&replicaSet=rs0"
	} */
	//mongoURI := fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, mongoPassword, mongoHost)
	//fmt.Println("mongoURI", mongoURI)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/?replicaSet=rs0")
	// Note the replicaSet=rs0 option is essential.
	// mongodb://admin:password123@localhost:27017/

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	return client.Database("health-tech")
}
