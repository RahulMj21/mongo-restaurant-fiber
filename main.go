package main

import (
	"log"

	"github.com/RahulMj21/mongo-restaurant-fiber/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routes.TestRoutes(app)

	err := app.Listen(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
