package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/controllers"
	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/database"
	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/middleware"
	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.Data(database.Client, "Products"), database.Data(database.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRouts(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
