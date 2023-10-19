package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MongoDB := os.Getenv("MONGODB_URL")

	clientOption := options.Client().ApplyURI(MongoDB)

	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		fmt.Println("Error is:", err)
		log.Fatal(err)
	}
	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	var collection *mongo.Collection = client.Database(dbName).Collection(collectionName)
	return collection
}
