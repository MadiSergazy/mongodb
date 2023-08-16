package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "gopkg.in/mgo.v2/bson"
	"go.mongodb.org/mongo-driver/bson"
)

// start-restaurant-struct
type Restaurant struct {
	Name         string
	RestaurantId string        `bson:"restaurant_id,omitempty"`
	Cuisine      string        `bson:"cuisine,omitempty"`
	Address      interface{}   `bson:"address,omitempty"`
	Borough      string        `bson:"borough,omitempty"`
	Grades       []interface{} `bson:"grades,omitempty"`
}

// end-restaurant-struct

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

	// begin insertMany
	coll := client.Database("sample_restaurants").Collection("restaurants")
	newRestaurants := []interface{}{
		Restaurant{Name: "Rule of Thirds", Cuisine: "Japanese"},
		Restaurant{Name: "Madame Vo", Cuisine: "Vietnamese"},
	}

	result, err := coll.InsertMany(context.TODO(), newRestaurants)
	if err != nil {
		panic(err)
	}
	// end insertMany

	// When you run this file, it should print:
	// 2 documents inserted with IDs: ObjectID("..."), ObjectID("...")
	fmt.Printf("%d documents inserted with IDs:\n", len(result.InsertedIDs))
	for _, id := range result.InsertedIDs {
		fmt.Printf("\t%s\n", id)
	}

	// 	In your code, the filter variable is constructed using bson.M{"cuisine": "Japanese"}, which creates a filter using a map representation.
	// 	This works well with the coll.Find() function because the map can be easily interpreted as a filter.

	// On the other hand, the filter2 variable is constructed using bson.D{{"cuisine", "Japanese"}},
	// which creates a filter using an array representation of key-value pairs.
	// The bson.D format is used when you need to maintain a specific order of the elements in the filter.
	// ^ However, in your code, you are using the old gopkg.in/mgo.v2/bson package to construct the filter,
	//  ^ which doesn't align with the new go.mongodb.org/mongo-driver/mongo package you are using for other operations.
	filter2 := bson.D{{"cuisine", "Japanese"}}
	//! JUST USE 	"go.mongodb.org/mongo-driver/bson"
	// filter := bson.M{"cuisine": "Japanese"}

	cursor, err := coll.Find(context.TODO(), filter2)
	if err != nil {
		panic(err)
	}
	// end find

	for cursor.Next(context.TODO()) {
		var result Restaurant
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
