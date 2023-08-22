package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"

	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/models"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "invalid Index " + http.StatusText(http.StatusNotFound)})
			c.Abort()
			return
		}

		userHexID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		}
		var address models.Address
		address.Address_id = primitive.NewObjectID()

		if err = c.BindJSON(&address); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, http.StatusText(http.StatusNotAcceptable))
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		match_filter := bson.D{primitive.E{Key: "_id", Value: userHexID}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id}}"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointcursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			c.Abort()
			return
		}

		var addressInfo []bson.M
		if err = pointcursor.All(ctx, addressInfo); err != nil {
			panic(err)
		}

		var size int32

		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)

		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: address}}}}
			_, err = userCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			c.IndentedJSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}
		// ctx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("user_id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid"})
			c.Abort()
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}
		var editaddress models.Address
		if err = c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set",
			Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House},
				{Key: "address.0.street_name", Value: editaddress.Street},
				{Key: "address.0.city_name", Value: editaddress.City},
				{Key: "adress.0.pin_code", Value: editaddress.Pincode}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
			return
		}
		c.IndentedJSON(200, "Successfully EditHomeAddress!")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("user_id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid"})
			c.Abort()
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}
		var editaddress models.Address
		if err = c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set",
			Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editaddress.House},
				{Key: "address.1.street_name", Value: editaddress.Street},
				{Key: "address.1.city_name", Value: editaddress.City},
				{Key: "adress.1.pin_code", Value: editaddress.Pincode}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		// defer cancel()
		// ctx.Done()
		c.IndentedJSON(200, "Successfully EditWorkAddress!")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "Wromg")
			return
		}
		// defer cancel()
		// ctx.Done()
		c.IndentedJSON(200, "Successfully Deleted!")
	}
}
