package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/controllers"
	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/database"
	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/middleware"
	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/routes"
)

func main() {
	port := os.GetEnv("PORT")
	if port == "" {
		port = "8080"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRouts(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddCart())
	router.GET("/removeitem", app.removeItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
