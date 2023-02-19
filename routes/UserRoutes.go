package routes

import (
	"github.com/RahulMj21/mongo-restaurant-fiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(api fiber.Router) {
	api.Get("/users", controllers.GetUsers)
	api.Get("/users/:id", controllers.GetUser)
	api.Post("/signup", controllers.Signup)
	api.Post("/signin", controllers.Signin)
}
