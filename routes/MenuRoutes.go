package routes

import (
	"github.com/RahulMj21/mongo-restaurant-fiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func MenuRoutes(app *fiber.App) {
	app.Get("/menus", controllers.GetMenus)
}
