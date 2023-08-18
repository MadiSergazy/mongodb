package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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

		address := make([]models.Address, 0)
		userHexID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: userHexID}}
		// update := bson.M{"$push": bson.M{"address": address}}
		update := bson.D{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: address}}}

		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			c.Abort()
			return
		}
		// ctx.Done()
		c.IndentedJSON(http.StatusOK, "Successfully deleted")
	}
}

func EditHomeAddress() gin.HandlerFunc {

}

func EditWorkAddress() gin.HandlerFunc {

}

func DeleteAddress() gin.HandlerFunc {

}
