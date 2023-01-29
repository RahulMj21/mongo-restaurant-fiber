package routes

import "github.com/gofiber/fiber/v2"

func TestRoutes(api fiber.Router) {
	api.Get("/health-check", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(&fiber.Map{
			"success": true,
			"message": "Bolo Radhe Shyam Ki Jayyyy.....",
		})
	})
}
