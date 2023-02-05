package main

import (
	"log"

	"github.com/RahulMj21/mongo-restaurant-fiber/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	api := app.Group("/api/v1")

	routes.TestRoutes(api)
	routes.FoodRoutes(api)
	routes.MenuRoutes(api)
	routes.OrderRoutes(api)

	err := app.Listen(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
