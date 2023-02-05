package routes

import (
	"github.com/RahulMj21/mongo-restaurant-fiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func OrderRoutes(api fiber.Router) {
	api.Get("/orders", controllers.GetOrders)
	api.Get("/orders/:id", controllers.GetOrder)
	api.Post("/orders", controllers.CreateOrder)
	api.Patch("/orders/:id", controllers.UpdateFood)
}
