package main

import (
	"log"

	"github.com/RahulMj21/mongo-restaurant-fiber/database"
	"github.com/RahulMj21/mongo-restaurant-fiber/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDB()

	app := fiber.New()

	routes.TestRoutes(app)
	routes.FoodRoutes(app)

	err := app.Listen(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
