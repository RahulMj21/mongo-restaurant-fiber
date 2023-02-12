package routes

import (
	"github.com/RahulMj21/mongo-restaurant-fiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func TableRoutes(api fiber.Router) {
	api.Get("/tables", controllers.GetTables)
	api.Get("/tables/:id", controllers.GetTable)
	api.Post("/tables", controllers.CreateTable)
	api.Patch("/tables/:id", controllers.UpdateTable)
}
