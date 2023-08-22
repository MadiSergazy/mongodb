package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"

	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/database"
	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/models"
)

const (
	emptyProduct       = "Product id is empty"
	emptyUserID        = "Product id is empty"
	notExistingProduct = "Product id is not existing"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection *mongo.Collection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println(emptyProduct)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println(emptyUserID)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New(emptyUserID))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(notExistingProduct+": ", err)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is not existing"))
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			// log.Println("Product id is not existing: ", err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Sucsessfully added product")

	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println("Product id is not existing: ", err)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is not existing"))
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			// log.Println("Product id is not existing: ", err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Sucsessfully removed item from cart")

	}
}

func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("content-type", "application/json")
			c.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			c.Abort()
			return

		}

		userHexID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		var filledcart models.User
		userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userHexID}}).Decode(&filledcart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		fillter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userHexID}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id}}"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

		pointcursor, err := app.userCollection.Aggregate(ctx, mongo.Pipeline{fillter_match, unwind, grouping})
		if err != nil {
			log.Print(err)
		}

		var listing []bson.M

		for pointcursor.Next(ctx) {
			var list bson.M

			err = pointcursor.Decode(&list)
			if err != nil {
				log.Print(err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			listing = append(listing, list)

		}
		defer pointcursor.Close(ctx)
		for _, json := range listing {
			c.IndentedJSON(http.StatusOK, json["total"])
			c.IndentedJSON(http.StatusOK, filledcart.UserCart)
		}
		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		if userQueryID == "" {
			log.Panic("User id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			// log.Println("Product id is not existing: ", err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Sucsessfully placed order")

	}

}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println("Product id is not existing: ", err)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is not existing"))
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			// log.Println("Product id is not existing: ", err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Sucsessfully placed order")
	}
}
