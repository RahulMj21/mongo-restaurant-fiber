package routes

import (
	"github.com/RahulMj21/mongo-restaurant-fiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func MenuRoutes(app *fiber.App) {
	app.Get("/menus", controllers.GetMenus)
	app.Get("/menus/:id", controllers.GetMenu)
	app.Post("/menus", controllers.CreateMenu)
	app.Patch("/menus/:id", controllers.UpdateMenu)
}
