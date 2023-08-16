package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// begin updatemany
	coll := client.Database("sample_airbnb").Collection("listingsAndReviews")
	filter := bson.D{{"address.market", "Sydney"}}
	update := bson.D{{"$mul", bson.D{{"price", 1.15}}}}

	result, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	// end updatemany

	// When you run this file for the first time, it should print:
	// Number of documents replaced: 609
	fmt.Printf("Documents updated: %v\n", result.ModifiedCount)

	cursor, err := coll.Find(context.TODO(), filter, options.Find().SetLimit(1))
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO()) // Make sure to close the cursor when you're done

	for cursor.Next(context.TODO()) {
		var result interface{}
		if err = cursor.Decode(&result); err != nil {
			panic(err)
		}

		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", output)

	}
}

type Restaurant struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string
	RestaurantId string `bson:"restaurant_id"`
	Cuisine      string
	Address      interface{}
	Borough      string
	Grades       interface{}
}
