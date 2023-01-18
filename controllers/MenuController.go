package controllers

import (
	"github.com/RahulMj21/mongo-restaurant-fiber/database"
	"github.com/RahulMj21/mongo-restaurant-fiber/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

var MenuCollection = database.OpenCollection(database.Client, "menu")

func GetMenus(c *fiber.Ctx) error {
	var menus []models.Menu

	cursor, err := MenuCollection.Find(c.Context(), bson.D{{}})
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	if err := cursor.All(c.Context(), &menus); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"status": "success",
		"data":   menus,
	})
}
