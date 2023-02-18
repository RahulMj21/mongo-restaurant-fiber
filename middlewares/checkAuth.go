package middlewares

import "github.com/gofiber/fiber/v2"

func Authorization(c *fiber.Ctx) {
	c.Next()
}
