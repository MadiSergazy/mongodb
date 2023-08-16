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

//* representation in mongoDB

//^ "_id":{"$oid":"5eb3d668b31de5d588f4292b"},
// "address":{"building":"7114","coord":
// [{"$numberDouble":"-73.9068506"},{"$numberDouble":"40.6199034"}],
// "street":"Avenue U","zipcode":"11234"},
// "borough":"Brooklyn",
// "cuisine":"Delicatessen",
// "grades":[{"date":{"$date":{"$numberLong":"1401321600000"}},
// "grade":"A","score":{"$numberInt":"10"}},
// {"date":{"$date":{"$numberLong":"1389657600000"}}],
// "name":"Wilken'S Fine Food",

// ^ "restaurant_id":"40356483"}
// start-restaurant-struct
type Restaurant struct {
	ID           primitive.ObjectID `bson:"_id"` // MongoDB uses _id as the default primary key field
	Name         string
	RestaurantId string `bson:"restaurant_id"`
	Cuisine      string
	Address      interface{}
	Borough      string
	Grades       []interface{}
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

	// begin findOne
	coll := client.Database("sample_restaurants").Collection("restaurants")
	filter := bson.D{{"name", "Bagels N Buns"}}
	// filter2 := bson.M{"name": "Bagels N Buns"}
	// * If you don't rely on the order of fields, you can use bson.M for simplicity. If order matters or you plan to extend the query with more complex operators, bson.D might be a better choice.
	var result Restaurant
	err = coll.FindOne(context.TODO(), filter).Decode(&result) //returning the first document matched

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return
		}
		panic(err)
	}
	// end findOne

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", output)
}
