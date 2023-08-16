package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoURI = "mongodb://localhost:27017"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	demoDB := client.Database("demo")
	err = demoDB.CreateCollection(ctx, "cats")
	if err != nil {
		log.Fatal(err)
	}
	catsCollection := demoDB.Collection("cats")
	defer catsCollection.Drop(ctx)
	result, err := catsCollection.InsertOne(ctx, bson.D{ //order matter
		{Key: "name", Value: "Mocha"},
		{Key: "breed", Value: "Turkish Van"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Result:", result)
	manyResult, err := catsCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{Key: "name", Value: "Latte"},
			{Key: "breed", Value: "Maine Coon"},
		},
		bson.D{
			{Key: "name", Value: "Trouble"},
			{Key: "breed", Value: "Domestic Shorthair"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Result:", manyResult)
	cursor, err := catsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var cats []bson.M                             //* order did't matter
	if err = cursor.All(ctx, &cats); err != nil { // ? not efficient for large data
		log.Fatal(err)
	}
	fmt.Println("Cats:", cats)
	cursor, err = catsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var kitty bson.M
		if err = cursor.Decode(&kitty); err != nil { //^efficient for large data
			log.Fatal(err)
		}
		fmt.Println("Kitty:", kitty)
	}
	var cat bson.M
	if err = catsCollection.FindOne(ctx, bson.M{}).Decode(&cat); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Cat:", cat)
	filter := bson.M{"breed": "Turkish Van"}
	fCursor, err := catsCollection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	var vans []bson.M
	if err = fCursor.All(ctx, &vans); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Turkish Vans:", vans[1])
}
