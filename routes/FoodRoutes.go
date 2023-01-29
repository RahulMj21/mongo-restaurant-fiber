package routes

import (
	"github.com/RahulMj21/mongo-restaurant-fiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func FoodRoutes(api fiber.Router) {
	api.Get("/foods", controllers.GetFoods)
	api.Get("/food/:id", controllers.GetFood)
	api.Post("/foods", controllers.CreateFood)
	api.Patch("/foods/:id", controllers.UpdateFood)
}
