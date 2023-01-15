package routes

import "github.com/gofiber/fiber/v2"

func TestRoutes(app *fiber.App) {
	app.Get("/health-check", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(&fiber.Map{
			"success": true,
			"message": "Bolo Radhe Shyam Ki Jayyyy.....",
		})
	})
}
