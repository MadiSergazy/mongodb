package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"

	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/database"
	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/models"
	generate "github.com/MadiSergazy/mongodb/GolangEcommerceProject/tokens"
)

var userCollection *mongo.Collection = database.Data(database.Client, "Users")
var productCollection *mongo.Collection = database.Data(database.Client, "Products")
var Validate = validator.New()

func HashPasssword(password string) string {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}

	return string(passwordHash)
}

func VerifyPasssword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(givenPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login or Password is not correct"
		valid = false
	}
	return valid, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil { //parse the incoming JSON request body and populate the user object with the data from the request
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		if validationErr := Validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": validationErr})
			return

		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "user already exists"})

		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "this phone number already exists"})
			return
		}
		password := HashPasssword(*user.Password)
		user.Password = &password
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)

		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		_, insertter := userCollection.InsertOne(ctx, user)

		if insertter != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "User did't created"})
			return
		}

		c.JSON(http.StatusCreated, "Sucsessfully sign in!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var context, cancel = context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		var founduser models.User
		if err := userCollection.FindOne(context, bson.M{"email": user.Email}).Decode(&founduser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Password incorrect"})
			return
		}
		PasswoedIsValid, msg := VerifyPasssword(*user.Password, *founduser.Password)
		if !PasswoedIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": msg})
			return
		}

		token, refreshtoken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)
		generate.UpdateAllTokens(token, refreshtoken, founduser.User_ID)
		c.JSON(http.StatusFound, founduser)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		var products models.Product

		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		products.Product_ID = primitive.NewObjectID()

		_, err := productCollection.InsertOne(ctx, products)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "successfully added")

	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		cursor, err := productCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "somesing went wrong pleace try after some time")
		}

		for cursor.Next(ctx) {
			var product models.Product
			if err = cursor.Decode(&product); err != nil {
				log.Println("SearchProduct decode err: ", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			productList = append(productList, product)
		}
		defer cursor.Close(context.TODO())

		if err := cursor.Err(); err != nil {
			log.Println("SearchProduct cursor err: ", err)
			c.IndentedJSON(http.StatusBadRequest, "invalid value")
			return
		}
		c.IndentedJSON(http.StatusOK, productList)
	}
}

func SearchProductByQuery() gin.HandlerFunc {

	return func(c *gin.Context) {
		var searchProducts []models.Product

		queryParam := c.Query("name")

		if queryParam == "" {
			log.Println("Query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error ": "Invalid search index"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		cursor, err := productCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		for cursor.Next(ctx) {

			var searchProduct models.Product
			if err = cursor.Decode(&searchProduct); err != nil {
				log.Println("SearchProduct decode err: ", err)
				c.IndentedJSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
				return
			}

			searchProducts = append(searchProducts, searchProduct)
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println("Error in SearchProductByQuery: ", err)
			c.IndentedJSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		c.IndentedJSON(http.StatusOK, searchProducts)
	}

}
