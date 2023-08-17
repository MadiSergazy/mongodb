package main

import (
	"mado/controllers"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController(GetSession())
	r.GET("")
	r.POST("")
	r.DELETE("")
}

func GetSession() *mongo.Session {

}
