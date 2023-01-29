package main

import (
	"log"

	"github.com/RahulMj21/mongo-restaurant-fiber/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	api := app.Group("/api/v1")

	routes.TestRoutes(api)
	routes.FoodRoutes(api)
	routes.MenuRoutes(api)

	err := app.Listen(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
