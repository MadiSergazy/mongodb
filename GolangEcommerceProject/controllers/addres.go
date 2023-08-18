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
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "invalid Index " + http.StatusText(http.StatusNotFound)})
			c.Abort()
			return
		}

		address := make([]models.Address, 0)
		user_hex_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: user_hex_id}}
		update := bson.D{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: address}}}

		_, err := userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Successfully deleted")
	}
}

func EditHomeAddress() gin.HandlerFunc {

}

func EditWorkAddress() gin.HandlerFunc {

}

func DeleteAddress() gin.HandlerFunc {

}
