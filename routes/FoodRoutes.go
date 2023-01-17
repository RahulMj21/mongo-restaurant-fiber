package routes

import (
	"github.com/RahulMj21/mongo-restaurant-fiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func FoodRoutes(app *fiber.App) {
	app.Get("/food/:id", controllers.GetFood)
	app.Post("/foods", controllers.CreateFood)
}
