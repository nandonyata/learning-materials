package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	router := gin.Default()
	// uri := "mongodb://admin:psw@localhost:27017/?authSource=admin" old
	uri := "mongodb://admin:psw@mongo_db:27017/?authSource=admin" // new, use "mongo_db" (container's name) instead of localhost

	fmt.Println(">>> AWOWIWI", uri)

	// Set up a client and connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to create Mongo client: %v", err)
	}

	// Connect to the MongoDB server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Get a handle for the "test" database and the "people" collection
	collection := client.Database("test-db").Collection("random-id-generator")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong 2",
		})
	})

	router.GET("/", func(c *gin.Context) {
		newPerson := bson.D{
			{"name", fmt.Sprintf("name-%v", time.Now().Unix())},
		}

		insertResult, err := collection.InsertOne(c, newPerson)
		if err != nil {
			log.Fatalf("Failed to insert document: %v", err)
		}

		fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)

		cursor, err := collection.Find(c, bson.D{})
		if err != nil {
			log.Fatalf("Failed to find document: %v", err)
		}
		defer cursor.Close(c)

		var results []bson.D
		for cursor.Next(c) {
			var document bson.D
			if err := cursor.Decode(&document); err != nil {
				log.Fatalf("Failed to decode document: %v", err)
			}
			results = append(results, document)
		}

		// Check for any cursor errors
		if err := cursor.Err(); err != nil {
			log.Fatalf("Cursor iteration error: %v", err)
		}

		c.JSON(200, gin.H{
			"message": "new person is inserted",
			"results": results,
		})
	})

	router.Run(":9002") // listen and serve on 0.0.0.0:9002
}
