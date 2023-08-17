package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"

	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/models"
)

func HashPasssword(password string) string {
	return ""
}

func VerifyPasssword(userPassword string, givenPassword string) (bool, string) {
	return false, ""
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
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.E})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "user already exists"})

		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

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
		user.UserCart = moke([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		_, insertter := UserCollection.InserOne(ctx, user)

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

		if err = UserCollection.FindOne(context, bson.M{"email": user.Email}).Decode(&founduser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Password incorrect"})
			return
		}
		PasswoedIsValid, msg := VerifyPasssword(*User.Password, *founduser.Password)
		if !PasswoedIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": msg})
			return
		}

		token, refreshtoken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *foundUser.Last_Name, *foundUser.User_ID)
		generate.UpdateAllTokens(token, refreshtoken, *founduser.User_ID)
		c.JSON(http.StatusFound, founduser)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {

}

func SearchProductByQuery() gin.HandlerFunc {
}
