package main

import (
	"context"
	"fmt"
	"log"

	"health-tech/lib/db"
	clog "health-tech/logger"
	"health-tech/user/services"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {

	userService := services.UserService{
		DB:     db.ConnectToDB(),
		Logger: *clog.MustNew("user", ""),
	}

	// Show database tables and rows
	fmt.Println("Database tables and rows:")
	fmt.Println("=========================")
	tableNames, err := userService.DB.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal("Failed to list collection names:", err)
	}
	for _, table := range tableNames {
		fmt.Println(table)
	}
	fmt.Println("=========================")
	// Get all rows from each table
	for _, table := range tableNames {
		fmt.Println(table)
		rows, err := userService.DB.Collection(table).Find(context.TODO(), bson.M{})
		if err != nil {
			log.Fatal("Failed to list collection names:", err)
		}
		for rows.Next(context.TODO()) {
			var row bson.M
			if err := rows.Decode(&row); err != nil {
				log.Fatal("Failed to decode row:", err)
			}
			fmt.Println(row)
		}
	}

}
