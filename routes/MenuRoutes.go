package routes

import (
	"github.com/RahulMj21/mongo-restaurant-fiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func MenuRoutes(api fiber.Router) {
	api.Get("/menus", controllers.GetMenus)
	api.Get("/menus/:id", controllers.GetMenu)
	api.Post("/menus", controllers.CreateMenu)
	api.Patch("/menus/:id", controllers.UpdateMenu)
}
